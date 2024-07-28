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

	if err := xml.Unmarshal(byteValue, &m); err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	m.Xmlns = m.XMLName.Space

	return
}

func (m *Manifest) GenerateModeDefault(output string) {
	os.MkdirAll(output, os.ModePerm)

}

func (m *Manifest) GenerateModeTypes(output string) {

	os.MkdirAll(output, os.ModePerm)
	for i, t := range m.Types {
		partManifest := Manifest{
			XMLName: m.XMLName,
			Xmlns:   m.Xmlns,
			Types:   []Types{t},
			Version: m.Version,
		}

		write(partManifest, output, i)
	}

}

func (m *Manifest) GenerateModeFileSize(output string, n int) {
	allTypes := m.Types

	componentsPerFile := int(math.Ceil(float64(len(allTypes)) / float64(n)))

	os.MkdirAll(output, os.ModePerm)
	for i := 0; i < n; i++ {
		startIdx := i * componentsPerFile
		endIdx := startIdx + componentsPerFile
		if endIdx > len(allTypes) {
			endIdx = len(allTypes)
		}
		TypesIntoXml := allTypes[startIdx:endIdx]

		if len(TypesIntoXml) == 0 {
			break
		}

		partManifest := Manifest{
			XMLName: m.XMLName,
			Xmlns:   m.Xmlns,
			Types:   TypesIntoXml,
			Version: m.Version,
		}

		write(partManifest, output, i)
	}
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
