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

func GenerateOutputDirectory(output string) (err error) {

	// 出力先ディレクトリの生成
	// ディレクトリが存在しない場合のみ作成
	err = os.MkdirAll(output, os.ModePerm)
	return

}

func (m *Manifest) GenerateXML(output string, mode string, n int) (err error) {

	// 1ファイルに含まれるコンポーネント数の取得
	componentsPerFile := m.calcComponentsPerFile(mode, n)

	// XML書き込み
	if len(m.Types) <= componentsPerFile {
		// コンポーネント数が1ファイルに含まれるコンポーネント数の取得以下のときはそのまま書き込む
		filename := generateFilename(output, nil)
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
		filename := generateFilename(output, &filenumber)

		err = partManifest.write(filename)
		if err != nil {
			return
		}
	}

	return

}

func (m *Manifest) GenerateXMLModeTypes(output string) (err error) {

	// Typesごとにpackage.xmlを分割する
	for i, t := range m.Types {
		i += 1 // ファイル番号

		partManifest := m.generatePartManifest([]Types{t})
		filename := generateFilename(output, &i)

		err = partManifest.write(filename)
		if err != nil {
			return
		}
	}

	return

}

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

func generateFilename(output string, filenumber *int) (filename string) {

	// ファイル番号が指定された場合は連番でファイルを生成する
	if filenumber == nil {
		filename = filepath.Join(output, Filename)
	} else {
		filename = filepath.Join(output, fmt.Sprintf(FilenameWithNumber, *filenumber))
	}

	return
}

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
