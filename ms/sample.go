package ms

import (
	"fmt"
	"strings"
)

func GenerateLargePackageXML(output string) {
	// メタデータタイプとメンバーを任意に指定
	metadataTypes := []string{"ApexClass", "CustomObject", "CustomField", "Workflow", "ValidationRule"}
	numMembers := MemberLimit + 1 // 各メタデータタイプのメンバー数を10001に設定

	var m Manifest
	m.Version = "61.0"

	for _, metadataType := range metadataTypes {
		var t Types
		t.Name = metadataType
		for i := 1; i <= numMembers; i++ {
			member := fmt.Sprintf("%s%d", strings.ToLower(metadataType), i)
			t.Members = append(t.Members, member)
		}
		m.Types = append(m.Types, t)
	}

	// XMLファイルを生成
	write(m, output, 1)
}
