package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
)

type Package struct {
	XMLName xml.Name `xml:"Package"`
	Xmlns   string   `xml:"xmln,attr"`
	Types   []Type   `xml:"types"`
	Version string   `xml:"version"`
}

type Type struct {
	Members []string `xml:"members"`
	Name    string   `xml:"name"`
}

func main() {

	input, output, num := recieveArgs()
	fmt.Println(*input, *output, *num)
	pkg := readXML(*input)
	splitByNumberOfFiles(pkg, *output, *num)

}

func recieveArgs() (inputManifest *string, outputDirectory *string, numFiles *int) {
	inputManifest = flag.String("input", "", "Path to the input XML file")
	outputDirectory = flag.String("output", "", "Path to the output directory")
	numFiles = flag.Int("numFiles", 1, "Number of files to split into")
	flag.Parse()

	if *inputManifest == "" || *outputDirectory == "" {
		flag.Usage()
		log.Fatal("Both input and output parameters are required")
	}

	if *numFiles == 1 {
		*numFiles = 2
	}

	return
}

func readXML(inputFile string) (pkg Package) {
	xmlFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening file %v", err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	if err := xml.Unmarshal(byteValue, &pkg); err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	return
}

func splitByType(pkg Package, outputDir string) {
	os.MkdirAll(outputDir, os.ModePerm)

	for i, t := range pkg.Types {
		partPackage := Package{
			XMLName: pkg.XMLName,
			Xmlns:   pkg.XMLName.Space,
			Types:   []Type{t},
			Version: pkg.Version,
		}

		resourceName := partPackage.Types[len(partPackage.Types)-1].Name

		outputFile := filepath.Join(outputDir, fmt.Sprintf("package_%03d_%v.xml", i+1, resourceName))
		output, err := xml.MarshalIndent(partPackage, "", "    ")
		if err != nil {
			log.Fatalf("Error wrighting file: %v", err)
		}

		err = os.WriteFile(outputFile, append([]byte(xml.Header), output...), 0644)

		fmt.Printf("Generated file: %s\n", outputFile)
	}
}

func splitByNumberOfFiles(pkg Package, outputDir string, numFiles int) {

	allComponents := pkg.Types
	componentsPerFile := int(math.Ceil(float64(len(allComponents)) / float64(numFiles)))

	os.MkdirAll(outputDir, os.ModePerm)

	for i := 0; i < numFiles; i++ {
		startIdx := i * componentsPerFile
		endIdx := startIdx + componentsPerFile
		if endIdx > len(allComponents) {
			endIdx = len(allComponents)
		}
		fileComponents := allComponents[startIdx:endIdx]

		if len(fileComponents) == 0 {
			break
		}

		partPackage := Package{
			XMLName: xml.Name{Local: "Package"},
			Types:   fileComponents,
			Version: pkg.Version,
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("package_%02d.xml", i+1))
		output, err := xml.MarshalIndent(partPackage, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling XML: %v", err)
		}

		err = os.WriteFile(outputFile, append([]byte(xml.Header), output...), 0644)
		if err != nil {
			log.Fatalf("Error writing file: %v", err)
		}

		fmt.Printf("Generated file: %s\n", outputFile)
	}

}
