package ms

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

// Package.xmlを格納
type Manifest struct {
	XMLName xml.Name `xml:"Package"`
	Xmlns   string   `xml:"xmlns,attr"`
	Types   []Types  `xml:"types"`
	Version string   `xml:"version"`
}

type Types struct {
	Members []string `xml:"members"`
	Name    string   `xml:"name"`
}

// package.xmlの読み込み
func ReadXML(input string) (m Manifest, err error) {

	xmlFile, err := os.Open(input)
	if err != nil {
		return
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(byteValue, &m); err != nil {
		return
	}

	m.Xmlns = m.XMLName.Space

	if m.XMLName.Local != "Package" || len(m.Types) == 0 {
		err = fmt.Errorf("有効なpackage.xmlではないか、typesが含まれていません")
		return
	}

	for _, t := range m.Types {
		if len(t.Members) == 0 {
			// membersを持たないtypesはSplitTypesで無言に消えてしまうため、ここで検出する
			err = fmt.Errorf("type %q にmembersが含まれていません", t.Name)
			return
		}
	}

	return

}

// 出力先ディレクトリの生成
func GenerateOutputDirectory(output string) (err error) {

	// ディレクトリが存在しない場合のみ作成
	err = os.MkdirAll(output, os.ModePerm)
	return

}

// 本ツールが生成するファイル名(package.xml または NNN_package.xml)と厳密に一致するかどうか
// "*package.xml"のような緩いglobだと無関係なファイル(例: backup_package.xml)まで対象にしてしまうため、
// 生成規則(constant.goのFilename/FilenameWithNumber)そのものに厳密一致させる
var numberedFilenamePattern = regexp.MustCompile(`^\d+_` + regexp.QuoteMeta(Filename) + `$`)

func isGeneratedFilename(name string) bool {
	return name == Filename || numberedFilenamePattern.MatchString(name)
}

// 出力先ディレクトリ内の、今回の実行で書き込まなかった残存生成物を削除する
// 「先に消してから書く」のではなく「書き終わった後に不要なものだけ消す」ことで、
// 読み込み/書き込みが失敗した場合に既存の正常な出力を巻き込んで失うことを防ぐ
// keepには今回書き込んだファイルの絶対/相対パスを渡す(このパスは削除対象から除外される)
// ディレクトリが存在しない場合は何もしない
func CleanOutputDirectory(output string, keep []string) (err error) {

	entries, err := os.ReadDir(output)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}

	keepSet := make(map[string]struct{}, len(keep))
	for _, k := range keep {
		if abs, absErr := filepath.Abs(k); absErr == nil {
			keepSet[abs] = struct{}{}
		}
	}

	for _, e := range entries {
		if e.IsDir() || !isGeneratedFilename(e.Name()) {
			continue
		}

		path := filepath.Join(output, e.Name())

		if abs, absErr := filepath.Abs(path); absErr == nil {
			if _, ok := keepSet[abs]; ok {
				continue
			}
		}

		if err = os.Remove(path); err != nil {
			return
		}
		color.Yellow("Removed: %s\n", path)
	}

	return

}

// xmlファイルを生成(default, files)
// 書き込んだファイルのパス一覧を返す(呼び出し元がCleanOutputDirectoryの keep として利用する)
func (m *Manifest) GenerateXML(output string, mode string, n int) (written []string, err error) {

	if mode != ModeDefault && mode != ModeFiles {
		err = fmt.Errorf("不正なモードです: %q", mode)
		return
	}

	if n < 1 {
		err = fmt.Errorf("コンポーネント数またはファイル数は1以上を指定してください: %d", n)
		return
	}

	// 1ファイルに含まれるコンポーネント数の取得
	componentsPerFile := m.calcComponentsPerFile(mode, n)

	// XML書き込み
	if len(m.Types) <= componentsPerFile {
		// コンポーネント数が1ファイルに含まれるコンポーネント数の取得以下のときはそのまま書き込む
		// SplitTypesで1メンバーごとに分解されているため、書き込み前に同名typeをまとめ直す
		m.combineTypes()
		filename := generateFilename(output)
		if err = m.write(filename); err != nil {
			return
		}
		written = append(written, filename)
		return
	}

	for i := 0; i < len(m.Types); i += componentsPerFile {
		// 1ファイルに含まれるコンポーネント数ごとにファイル書き込み
		end := i + componentsPerFile
		if end > len(m.Types) {
			end = len(m.Types)
		}
		partManifest := m.generatePartManifest(m.Types[i:end])

		filenumber := int(math.Ceil(float64(end) / float64(componentsPerFile)))
		filename := generateFilenameWithNumber(output, filenumber)

		partManifest.combineTypes()
		if err = partManifest.write(filename); err != nil {
			return
		}
		written = append(written, filename)
	}

	return

}

// Typesごとにpackage.xmlを分割する
// 書き込んだファイルのパス一覧を返す(呼び出し元がCleanOutputDirectoryの keep として利用する)
func (m *Manifest) GenerateXMLModeTypes(output string) (written []string, err error) {

	for i, t := range m.Types {
		i += 1 // ファイル番号

		partManifest := m.generatePartManifest([]Types{t})
		filename := generateFilenameWithNumber(output, i)

		if err = partManifest.write(filename); err != nil {
			return
		}
		written = append(written, filename)
	}

	return

}

// typesを1コンポーネントごとに分割する
func (m *Manifest) SplitTypes() {

	tmp := m.Types
	m.Types = []Types{}

	for _, types := range tmp {
		for _, member := range types.Members {
			typeToAppend := Types{
				Members: []string{member},
				Name:    types.Name,
			}
			m.Types = append(m.Types, typeToAppend)
		}
	}

}

// コンポーネントをNameごとにTypesにまとめる
// mapの走査順は実行のたびに変わるため、Nameの初出順を別途保持して出力順を安定させる
func (m *Manifest) combineTypes() {

	typesMap := make(map[string]*Types)
	var order []string

	for _, t := range m.Types {
		if types, ok := typesMap[t.Name]; ok {
			types.Members = append(types.Members, t.Members...)
			continue
		}

		typesCopy := Types{Name: t.Name, Members: append([]string{}, t.Members...)}
		typesMap[t.Name] = &typesCopy
		order = append(order, t.Name)
	}

	m.Types = make([]Types, 0, len(order))
	for _, name := range order {
		m.Types = append(m.Types, *typesMap[name])
	}

}

// 1ファイルに書き込むコンポーネントの上限を取得する
// mode は GenerateXML で ModeDefault/ModeFiles のいずれかであることを検証済みの前提のため、
// それ以外はここでは定義されていないmodeとして扱う
func (m *Manifest) calcComponentsPerFile(mode string, n int) (componentsPerFile int) {

	switch mode {
	case ModeDefault:
		return n
	case ModeFiles:
		return int(math.Ceil(float64(len(m.Types)) / float64(n)))
	default:
		panic(fmt.Sprintf("定義されていないmodeです: %q", mode))
	}

}

// ファイルに書き込む構造体を生成する
func (m *Manifest) generatePartManifest(types []Types) (partManifest Manifest) {

	partManifest = Manifest{
		XMLName: m.XMLName,
		Xmlns:   m.Xmlns,
		Types:   types,
		Version: m.Version,
	}

	return

}

// ファイル名を生成
func generateFilename(output string) (filename string) {

	return filepath.Join(output, Filename)

}

// 番号付きのファイル名を生成
func generateFilenameWithNumber(output string, filenumber int) (filename string) {

	return filepath.Join(output, fmt.Sprintf(FilenameWithNumber, filenumber))

}

// xmlファイル書き込み処理
func (m *Manifest) write(filename string) (err error) {

	manifestXml, err := xml.MarshalIndent(*m, "", "    ")
	if err != nil {
		return
	}

	err = os.WriteFile(filename, append([]byte(xml.Header), manifestXml...), 0644)
	if err != nil {
		return
	}

	color.Green("Generated: %s\n", filename)
	return

}
