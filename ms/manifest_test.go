package ms

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test fixture %s: %v", path, err)
	}
	return path
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to exist, but got: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed, but stat returned: %v", path, err)
	}
}

func newTestManifest() Manifest {
	return Manifest{
		Xmlns:   "http://soap.sforce.com/2006/04/metadata",
		Version: "61.0",
		Types: []Types{
			{Name: "ApexClass", Members: []string{"Foo", "Bar", "Baz"}},
			{Name: "CustomObject", Members: []string{"Obj1"}},
		},
	}
}

func findType(t *testing.T, types []Types, name string) Types {
	t.Helper()

	for _, ty := range types {
		if ty.Name == name {
			return ty
		}
	}
	t.Fatalf("type %q not found in %+v", name, types)
	return Types{}
}

// ---------- ReadXML ----------

func TestReadXML(t *testing.T) {

	t.Run("valid package.xml", func(t *testing.T) {
		dir := t.TempDir()
		path := writeTestFile(t, dir, "package.xml", `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <members>Foo</members>
        <members>Bar</members>
        <name>ApexClass</name>
    </types>
    <version>61.0</version>
</Package>`)

		m, err := ReadXML(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.Xmlns != "http://soap.sforce.com/2006/04/metadata" {
			t.Errorf("Xmlns = %q, want the metadata namespace", m.Xmlns)
		}
		if len(m.Types) != 1 || m.Types[0].Name != "ApexClass" {
			t.Fatalf("unexpected Types: %+v", m.Types)
		}
		if want := []string{"Foo", "Bar"}; !slices.Equal(m.Types[0].Members, want) {
			t.Errorf("Members = %v, want %v", m.Types[0].Members, want)
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := ReadXML(filepath.Join(t.TempDir(), "missing.xml"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("malformed xml", func(t *testing.T) {
		dir := t.TempDir()
		path := writeTestFile(t, dir, "broken.xml", `<Package><types>`)

		if _, err := ReadXML(path); err == nil {
			t.Fatal("expected error for malformed xml, got nil")
		}
	})

	t.Run("root element is not Package", func(t *testing.T) {
		dir := t.TempDir()
		path := writeTestFile(t, dir, "notpackage.xml", `<?xml version="1.0"?><NotPackage/>`)

		if _, err := ReadXML(path); err == nil {
			t.Fatal("expected error for non-Package root, got nil")
		}
	})

	t.Run("empty types", func(t *testing.T) {
		dir := t.TempDir()
		path := writeTestFile(t, dir, "empty.xml", `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <version>61.0</version>
</Package>`)

		if _, err := ReadXML(path); err == nil {
			t.Fatal("expected error for empty types, got nil")
		}
	})

	t.Run("type with zero members", func(t *testing.T) {
		dir := t.TempDir()
		path := writeTestFile(t, dir, "nomembers.xml", `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">
    <types>
        <name>SomeType</name>
    </types>
    <version>61.0</version>
</Package>`)

		_, err := ReadXML(path)
		if err == nil {
			t.Fatal("expected error for type with zero members, got nil")
		}
		if !strings.Contains(err.Error(), "SomeType") {
			t.Errorf("error message %q does not mention the offending type name", err.Error())
		}
	})
}

// ---------- GenerateOutputDirectory ----------

func TestGenerateOutputDirectory(t *testing.T) {

	t.Run("creates a directory that does not exist yet", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "nested", "output")
		if err := GenerateOutputDirectory(dir); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			t.Fatalf("expected %s to be a directory", dir)
		}
	})

	t.Run("is a no-op when the directory already exists", func(t *testing.T) {
		dir := t.TempDir()
		if err := GenerateOutputDirectory(dir); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// ---------- isGeneratedFilename / CleanOutputDirectory ----------

func TestIsGeneratedFilename(t *testing.T) {

	tests := []struct {
		name string
		want bool
	}{
		{"package.xml", true},
		{"001_package.xml", true},
		{"0_package.xml", true},
		{"1000_package.xml", true}, // 3桁を超える連番にもマッチする必要がある
		{"backup_package.xml", false},
		{"package.xml.bak", false},
		{"_package.xml", false},    // 数字プレフィックスがない
		{"abc_package.xml", false}, // プレフィックスが数字でない
		{"Package.xml", false},     // 大文字小文字は区別する
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGeneratedFilename(tt.name); got != tt.want {
				t.Errorf("isGeneratedFilename(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestCleanOutputDirectory(t *testing.T) {

	t.Run("directory does not exist yet", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "does-not-exist")
		if err := CleanOutputDirectory(dir, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("removes generated files that are not in keep", func(t *testing.T) {
		dir := t.TempDir()
		stale1 := writeTestFile(t, dir, "001_package.xml", "x")
		stale2 := writeTestFile(t, dir, "002_package.xml", "x")
		keep := writeTestFile(t, dir, "package.xml", "x")

		if err := CleanOutputDirectory(dir, []string{keep}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertNotExists(t, stale1)
		assertNotExists(t, stale2)
		assertExists(t, keep)
	})

	t.Run("never touches files that do not match the generated filename pattern", func(t *testing.T) {
		// 過去の実装は "*package.xml" という緩いglobを使っており、
		// backup_package.xml のような無関係なファイルまで削除してしまっていた(finding #1の回帰確認)
		dir := t.TempDir()
		unrelated := writeTestFile(t, dir, "backup_package.xml", "x")
		writeTestFile(t, dir, "001_package.xml", "x") // keepに含まれないため削除される想定

		if err := CleanOutputDirectory(dir, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertExists(t, unrelated)
	})

	t.Run("never removes a directory even if its name matches the pattern", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "package.xml")
		if err := os.Mkdir(sub, 0755); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := CleanOutputDirectory(dir, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertExists(t, sub)
	})

	t.Run("protects the source file when input==output and it was just (re)written", func(t *testing.T) {
		// -clean と input==output の組み合わせでも、今回書き込んだファイル(keep)は保護される(finding #2/#3の回帰確認)
		dir := t.TempDir()
		self := writeTestFile(t, dir, "package.xml", "original-content")

		if err := CleanOutputDirectory(dir, []string{self}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertExists(t, self)
	})
}

// ---------- GenerateXML ----------

func TestManifest_GenerateXML(t *testing.T) {

	t.Run("n=0 is an error in default mode", func(t *testing.T) {
		m := newTestManifest()
		m.SplitTypes()
		if _, err := m.GenerateXML(t.TempDir(), ModeDefault, 0); err == nil {
			t.Fatal("expected error for n=0, got nil")
		}
	})

	t.Run("n=0 is an error in files mode", func(t *testing.T) {
		m := newTestManifest()
		m.SplitTypes()
		if _, err := m.GenerateXML(t.TempDir(), ModeFiles, 0); err == nil {
			t.Fatal("expected error for n=0, got nil")
		}
	})

	t.Run("an unknown mode is an error", func(t *testing.T) {
		m := newTestManifest()
		m.SplitTypes()
		if _, err := m.GenerateXML(t.TempDir(), "bogus", 10); err == nil {
			t.Fatal("expected error for unknown mode, got nil")
		}
	})

	t.Run("同名typeが1ファイルに収まる場合、1ブロックに統合される", func(t *testing.T) {
		// SplitTypesで1メンバーごとに分解した後、単一ファイルに収まるケースでcombineTypesが
		// 抜けていたためtypesブロックが重複生成されていたバグの回帰テスト
		dir := t.TempDir()
		m := newTestManifest()
		m.SplitTypes()

		written, err := m.GenerateXML(dir, ModeDefault, 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(written) != 1 {
			t.Fatalf("written = %v, want exactly 1 file", written)
		}

		got, err := ReadXML(written[0])
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}
		if len(got.Types) != 2 {
			t.Fatalf("Types count = %d, want 2 (merged by name), got %+v", len(got.Types), got.Types)
		}

		apex := findType(t, got.Types, "ApexClass")
		if want := []string{"Foo", "Bar", "Baz"}; !slices.Equal(apex.Members, want) {
			t.Errorf("ApexClass members = %v, want %v (in original order)", apex.Members, want)
		}
	})

	t.Run("複数ファイルに分割してもメンバーが欠損しない", func(t *testing.T) {
		dir := t.TempDir()
		m := newTestManifest()
		m.SplitTypes() // Foo, Bar, Baz, Obj1 の4エントリに分解される

		written, err := m.GenerateXML(dir, ModeDefault, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(written) != 4 {
			t.Fatalf("written = %v, want 4 files", written)
		}

		total := 0
		for _, f := range written {
			assertExists(t, f)
			got, err := ReadXML(f)
			if err != nil {
				t.Fatalf("failed to read %s: %v", f, err)
			}
			for _, ty := range got.Types {
				total += len(ty.Members)
			}
		}
		if total != 4 {
			t.Errorf("total members across files = %d, want 4", total)
		}
	})

	t.Run("filesモードは指定したファイル数になるよう分割する", func(t *testing.T) {
		dir := t.TempDir()
		m := newTestManifest()
		m.SplitTypes() // 4エントリ

		written, err := m.GenerateXML(dir, ModeFiles, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(written) != 2 {
			t.Fatalf("written = %v, want 2 files", written)
		}
	})

	t.Run("生成ファイルにXML宣言が付与される", func(t *testing.T) {
		dir := t.TempDir()
		m := newTestManifest()
		m.SplitTypes()

		written, err := m.GenerateXML(dir, ModeDefault, 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		content, err := os.ReadFile(written[0])
		if err != nil {
			t.Fatalf("failed to read generated file: %v", err)
		}
		if !strings.HasPrefix(string(content), xml.Header) {
			t.Errorf("generated file does not start with the XML header: %q", string(content[:min(len(content), 60)]))
		}
	})
}

// ---------- GenerateXMLModeTypes ----------

func TestManifest_GenerateXMLModeTypes(t *testing.T) {
	dir := t.TempDir()
	m := newTestManifest()

	written, err := m.GenerateXMLModeTypes(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		filepath.Join(dir, "001_package.xml"),
		filepath.Join(dir, "002_package.xml"),
	}
	if !slices.Equal(written, want) {
		t.Fatalf("written = %v, want %v", written, want)
	}

	got, err := ReadXML(written[0])
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}
	if len(got.Types) != 1 || got.Types[0].Name != "ApexClass" {
		t.Errorf("001_package.xml Types = %+v, want a single ApexClass block", got.Types)
	}
}

// ---------- SplitTypes ----------

func TestManifest_SplitTypes(t *testing.T) {

	t.Run("splits each member into its own Types entry", func(t *testing.T) {
		m := Manifest{Types: []Types{
			{Name: "ApexClass", Members: []string{"Foo", "Bar"}},
			{Name: "CustomObject", Members: []string{"Obj1"}},
		}}

		m.SplitTypes()

		want := []Types{
			{Name: "ApexClass", Members: []string{"Foo"}},
			{Name: "ApexClass", Members: []string{"Bar"}},
			{Name: "CustomObject", Members: []string{"Obj1"}},
		}
		if len(m.Types) != len(want) {
			t.Fatalf("SplitTypes() = %+v, want %+v", m.Types, want)
		}
		for i := range want {
			if m.Types[i].Name != want[i].Name || !slices.Equal(m.Types[i].Members, want[i].Members) {
				t.Errorf("Types[%d] = %+v, want %+v", i, m.Types[i], want[i])
			}
		}
	})

	t.Run("membersが空のtypeは無言で消える", func(t *testing.T) {
		// ReadXMLが事前にmembers空のtypeを拒否するため実運用では到達しないが、
		// SplitTypes単体の挙動として明示しておく(この前提が崩れるとfinding #5が再発する)
		m := Manifest{Types: []Types{
			{Name: "NoMembers", Members: nil},
			{Name: "ApexClass", Members: []string{"Foo"}},
		}}

		m.SplitTypes()

		if len(m.Types) != 1 || m.Types[0].Name != "ApexClass" {
			t.Fatalf("SplitTypes() = %+v, want only the ApexClass entry", m.Types)
		}
	})
}

// ---------- combineTypes ----------

func TestManifest_combineTypes_deterministicOrder(t *testing.T) {
	// mapの走査順に依存していた旧実装では実行のたびに<types>の並びが変わっていた(finding #バグ2の回帰確認)
	names := []string{"TypeA", "TypeB", "TypeC", "TypeD", "TypeE", "TypeF", "TypeG", "TypeH"}

	build := func() Manifest {
		m := Manifest{}
		for _, n := range names {
			m.Types = append(m.Types, Types{Name: n, Members: []string{"m1", "m2"}})
		}
		m.SplitTypes()
		return m
	}

	var firstOrder []string
	for i := 0; i < 20; i++ {
		m := build()
		m.combineTypes()

		var order []string
		for _, ty := range m.Types {
			order = append(order, ty.Name)
		}

		if firstOrder == nil {
			firstOrder = order
			continue
		}
		if !slices.Equal(order, firstOrder) {
			t.Fatalf("combineTypes order is not deterministic: run %d got %v, want %v", i, order, firstOrder)
		}
	}

	if !slices.Equal(firstOrder, names) {
		t.Errorf("combineTypes order = %v, want original insertion order %v", firstOrder, names)
	}
}

// ---------- calcComponentsPerFile ----------

func TestManifest_calcComponentsPerFile(t *testing.T) {

	tests := []struct {
		name  string
		types int
		mode  string
		n     int
		want  int
	}{
		{"default mode returns n as-is", 5, ModeDefault, 3, 3},
		{"files mode ceils the division", 10, ModeFiles, 3, 4},
		{"files mode with an exact division", 10, ModeFiles, 5, 2},
		{"files mode with n larger than the type count", 3, ModeFiles, 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manifest{Types: make([]Types, tt.types)}
			if got := m.calcComponentsPerFile(tt.mode, tt.n); got != tt.want {
				t.Errorf("calcComponentsPerFile(%q, %d) with %d types = %d, want %d", tt.mode, tt.n, tt.types, got, tt.want)
			}
		})
	}

	t.Run("panics on an unknown mode", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected a panic for an unknown mode, got none")
			}
		}()
		m := Manifest{Types: make([]Types, 1)}
		m.calcComponentsPerFile("bogus", 1)
	})
}

// ---------- filename helpers ----------

func TestGenerateFilename(t *testing.T) {
	got := generateFilename("out")
	want := filepath.Join("out", "package.xml")
	if got != want {
		t.Errorf("generateFilename() = %q, want %q", got, want)
	}
}

func TestGenerateFilenameWithNumber(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "001_package.xml"},
		{42, "042_package.xml"},
		{1000, "1000_package.xml"}, // %03dは最小桁数の指定なので4桁以上に伸びる
	}
	for _, tt := range tests {
		got := generateFilenameWithNumber("out", tt.n)
		want := filepath.Join("out", tt.want)
		if got != want {
			t.Errorf("generateFilenameWithNumber(%d) = %q, want %q", tt.n, got, want)
		}
	}
}
