package ms

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Package.xmlを格納
// 読み取り、書き込みともに利用する
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

func ReadXML(input string) (m Manifest, err error) {

	// package.xmlの読み込み
	xmlFile, err := os.Open(input)
	if err != nil {
		err = fmt.Errorf("error opening %v %w", input, err)
		return
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	if err = xml.Unmarshal(byteValue, &m); err != nil {
		err = fmt.Errorf("error unmarshalling xml %v", err)
		return
	}

	m.Xmlns = m.XMLName.Space

	return

}

func GenerateOutputDirectory(output string) (err error) {

	// 出力先ディレクトリの生成
	// ディレクトリが存在しない場合のみ作成
	err = os.MkdirAll(output, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("error making output directory: %v", err)
	}

	return

}

func (m *Manifest) GenerateXML(output string, mode string, n int) (err error) {

	// typesをコンポーネント毎に分割する
	m.splitTypes()

	// 1ファイルに含まれるコンポーネント数の取得
	componentsPerFile := m.calcComponentsPerFile(mode, n)

	// XML書き込み
	if len(m.Types) <= componentsPerFile {
		// コンポーネント数が上限以下のときはそのまま書き込む
		m.write(output, nil)

	} else {
		// 指定されたコンポーネント数以下のpackage.xmlを作成する
		typesToWrite := []Types{}
		for i, t := range m.Types {
			i += 1 // 処理済みコンポーネント数
			typesToWrite = append(typesToWrite, t)

			if len(typesToWrite) == componentsPerFile || i == len(m.Types) {
				// コンポーネント数が１ファイルの上限を超えたとき、もしくはループが最後のコンポーネントまで達したとき
				partManifest := m.generatePartManifest(typesToWrite)

				filenumber := int(math.Ceil(float64(i) / float64(componentsPerFile)))
				err = partManifest.write(output, &filenumber)
				if err != nil {
					return
				}

				typesToWrite = []Types{}
			}
		}
	}

	return

}

func (m *Manifest) GenerateXMLModeTypes(output string) (err error) {

	// Typesごとにpackage.xmlを分割する
	for i, t := range m.Types {
		i += 1
		partManifest := m.generatePartManifest([]Types{t})
		err = partManifest.write(output, &i)
		if err != nil {
			return
		}
	}

	return
}

func (m *Manifest) splitTypes() {

	tmp := m.Types
	m.Types = []Types{}

	for _, types := range tmp {
		name := types.Name
		for _, member := range types.Members {
			typeToAppend := Types{
				Members: []string{member},
				Name:    name,
			}
			m.Types = append(m.Types, typeToAppend)
		}
	}
}

func (m *Manifest) calcComponentsPerFile(mode string, n int) (componentsPerFile int) {

	// 1ファイルに書き込むコンポーネントの上限を取得する
	switch mode {
	case ModeDefault:
		return n
	case ModeFiles:
		return int(math.Ceil(float64(len(m.Types)) / float64(n)))
	default:
		return MemberLimit
	}

}

func (m *Manifest) generatePartManifest(types []Types) (partManifest Manifest) {

	// ファイルに書き込む構造体を生成する
	partManifest = Manifest{
		XMLName: m.XMLName,
		Xmlns:   m.Xmlns,
		Types:   types,
		Version: m.Version,
	}

	return

}

func (m *Manifest) write(output string, filenumber *int) (err error) {

	// XMLファイルの生成
	var filename string

	// ファイル番号が指定された場合は連番でファイルを生成する
	if filenumber == nil {
		filename = filepath.Join(output, Filename)
	} else {
		filename = filepath.Join(output, fmt.Sprintf(FilenameWithNumber, *filenumber))
	}

	manifestXml, err := xml.MarshalIndent(*m, "", "    ")
	if err != nil {
		err = fmt.Errorf("error marshalling xml: %v", err)
		return
	}

	err = os.WriteFile(filename, append([]byte(xml.Header), manifestXml...), 0644)
	if err != nil {
		err = fmt.Errorf("error writing file: %v", err)
		return
	}

	color.Green("Generated: %s\n", filename)
	return

}
