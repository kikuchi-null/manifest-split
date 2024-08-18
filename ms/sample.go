package ms

import (
	"fmt"
	"strings"
)

// sampleのpackage.xmlを生成
func GenerateLargePackageXML(output string) (err error) {

	m := Manifest{
		Xmlns:   SampleXmlns,
		Version: SampleVersion,
	}

	for _, metadataType := range GetTypes() {
		var t Types
		t.Name = metadataType
		for i := 1; i <= SampleNum; i++ {
			member := fmt.Sprintf("%s%05d", strings.ToLower(metadataType), i)
			t.Members = append(t.Members, member)
		}
		m.Types = append(m.Types, t)
	}

	// XMLファイルを生成
	filename := generateFilename(output)
	err = m.write(filename)

	return

}
