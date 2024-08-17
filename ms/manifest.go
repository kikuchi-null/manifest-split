package ms

import (
	"encoding/xml"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"path/filepath"
	"slices"

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

	byteValue, _ := io.ReadAll(xmlFile)

	if err = xml.Unmarshal(byteValue, &m); err != nil {
		return
	}

	m.Xmlns = m.XMLName.Space

	return

}

// 出力先ディレクトリの生成
func GenerateOutputDirectory(output string) (err error) {

	// ディレクトリが存在しない場合のみ作成
	err = os.MkdirAll(output, os.ModePerm)
	return

}

// xmlファイルを生成(default, files)
func (m *Manifest) GenerateXML(output string, mode string, n int) (err error) {

	// 1ファイルに含まれるコンポーネント数の取得
	componentsPerFile := m.calcComponentsPerFile(mode, n)

	// XML書き込み
	if len(m.Types) <= componentsPerFile {
		// コンポーネント数が1ファイルに含まれるコンポーネント数の取得以下のときはそのまま書き込む
		filename := generateFilename(output)
		err = m.write(filename)
		return
	}

	for i := 0; i <= len(m.Types); i += componentsPerFile {
		// 1ファイルに含まれるコンポーネント数ごとにファイル書き込み
		end := i + componentsPerFile
		if end > len(m.Types) {
			end = len(m.Types)
		}
		partManifest := m.generatePartManifest(m.Types[i:end])

		filenumber := int(math.Ceil(float64(end) / float64(componentsPerFile)))
		filename := generateFilenameWithNumber(output, filenumber)

		partManifest.combineTypes()
		err = partManifest.write(filename)
		if err != nil {
			return
		}
	}

	return

}

// Typesごとにpackage.xmlを分割する
func (m *Manifest) GenerateXMLModeTypes(output string) (err error) {

	for i, t := range m.Types {
		i += 1 // ファイル番号

		partManifest := m.generatePartManifest([]Types{t})
		filename := generateFilenameWithNumber(output, i)

		err = partManifest.write(filename)
		if err != nil {
			return
		}
	}

	return

}

// typesを1コンポーネントごとに分割する
func (m *Manifest) SplitTypes() {

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

// コンポーネントをNameごとにTypesにまとめる
func (m *Manifest) combineTypes() {

	typesMap := make(map[string]Types)

	for _, t := range m.Types {

		if types, ok := typesMap[t.Name]; ok {
			types.Members = append(types.Members, t.Members...)
			typesMap[t.Name] = types
			continue
		}

		typesMap[t.Name] = t
	}

	m.Types = []Types{}
	m.Types = slices.AppendSeq(m.Types, maps.Values(typesMap))

}

// 1ファイルに書き込むコンポーネントの上限を取得する
func (m *Manifest) calcComponentsPerFile(mode string, n int) (componentsPerFile int) {

	switch mode {
	case ModeDefault:
		return n
	case ModeFiles:
		return int(math.Ceil(float64(len(m.Types)) / float64(n)))
	default:
		return MemberLimit
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
