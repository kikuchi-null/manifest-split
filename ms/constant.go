package ms

// 各種モード
const (
	ModeDefault string = "default"
	ModeFiles   string = "files"
	ModeTypes   string = "types"
	ModeSample  string = "sample"
	MemberLimit int    = 10000 // 1ファイルの上限
)

func GetModes() (modes []string) {
	return []string{ModeDefault, ModeFiles, ModeTypes, ModeSample}
}

// ファイル名
const (
	FilenameWithNumber string = "%03d_package.xml"
	Filename           string = "package.xml"
)

// サンプル生成に利用
const (
	SampleXmlns   string = "http://soap.sforce.com/2006/04/metadata"
	SampleVersion string = "61.0"
	SampleNum     int    = 1000 //コンポーネント数
)

func GetTypes() (types []string) {
	types = []string{
		"ApexClass",
		"ApexTrigger",
		"CustomApplication",
		"CustomObject",
		"CustomField",
		"Profile",
		"Workflow",
		"ValidationRule",
	}
	return
}
