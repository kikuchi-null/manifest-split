package ms

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
)

// Package.xmlを格納
// 読み取り、書き込みともに利用する
type Manifest struct {
	XMLName xml.Name `xml:"Package"`
	Xmlns   string   `xml:"xmln,attr"`
	Types   []Types  `xml:"types"`
	Version string   `xml:"version"`
}

type Types struct {
	Members []string `xml:"members"`
	Name    string   `xml:"name"`
}

func ReadXML(input string) (m Manifest) {
	xmlFile, err := os.Open(input)
	if err != nil {
		log.Fatalf("Error opening file %v", err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	if err = xml.Unmarshal(byteValue, &m); err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	m.Xmlns = m.XMLName.Space

	return

}

func GenerateOutputDirectory(output string) {

	err := os.MkdirAll(output, os.ModePerm)
	log.Fatalf("Error making output directory: %v", err)

}

func (m *Manifest) GenerateModeDefault(output string, componentsPerFile int) {

	// Typesをコンポーネント毎に分割
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

	if len(m.Types) <= componentsPerFile {
		write(*m, output, 1)

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
				write(partManifest, output, filenumber)

				typesToWrite = []Types{}
			}
		}
	}

}

func (m *Manifest) GenerateModeTypes(output string) {

	for i, t := range m.Types {
		partManifest := m.generatePartManifest([]Types{t})
		write(partManifest, output, i)
	}

}

func (m *Manifest) GenerateModeFileSize(output string, n int) {

	// 1ファイルごとのTypes数
	componentsPerFile := int(math.Ceil(float64(len(m.Types)) / float64(n)))

	for i := 0; i < n; i++ {
		startIdx := i * componentsPerFile
		endIdx := startIdx + componentsPerFile
		if endIdx > len(m.Types) {
			endIdx = len(m.Types)
		}
		typesToWrite := m.Types[startIdx:endIdx]

		if len(typesToWrite) == 0 {
			break
		}

		partManifest := m.generatePartManifest(typesToWrite)
		write(partManifest, output, i)
	}

}

func (m *Manifest) generatePartManifest(types []Types) (partManifest Manifest) {

	partManifest = Manifest{
		XMLName: m.XMLName,
		Xmlns:   m.Xmlns,
		Types:   types,
		Version: m.Version,
	}

	return
}

func write(manifest Manifest, output string, filenumber int) {

	finename := filepath.Join(output, fmt.Sprintf("%03d_package.xml", filenumber))
	manifestXml, err := xml.MarshalIndent(manifest, "", "    ")

	if err != nil {
		log.Fatalf("Error marshalling XML: %v", err)
	}

	err = os.WriteFile(finename, append([]byte(xml.Header), manifestXml...), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Printf("Generated file: %s\n", finename)

}
