package ms

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLargePackageXML(t *testing.T) {
	dir := t.TempDir()

	written, err := GenerateLargePackageXML(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantPath := filepath.Join(dir, Filename)
	if len(written) != 1 || written[0] != wantPath {
		t.Fatalf("written = %v, want [%s]", written, wantPath)
	}

	// 生成物自体がReadXMLで読める(=有効なpackage.xmlである)ことも合わせて確認する
	m, err := ReadXML(written[0])
	if err != nil {
		t.Fatalf("failed to read generated sample: %v", err)
	}

	types := GetTypes()
	if len(m.Types) != len(types) {
		t.Fatalf("Types count = %d, want %d", len(m.Types), len(types))
	}

	for i, ty := range m.Types {
		if ty.Name != types[i] {
			t.Errorf("Types[%d].Name = %q, want %q", i, ty.Name, types[i])
		}
		if len(ty.Members) != SampleNum {
			t.Errorf("%s members = %d, want %d", ty.Name, len(ty.Members), SampleNum)
		}
		if len(ty.Members) > 0 {
			wantFirst := fmt.Sprintf("%s%05d", strings.ToLower(ty.Name), 1)
			if ty.Members[0] != wantFirst {
				t.Errorf("%s first member = %q, want %q", ty.Name, ty.Members[0], wantFirst)
			}
		}
	}
}
