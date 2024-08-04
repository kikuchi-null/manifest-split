package ms

// 各種モード
const (
	ModeDefault string = "default"
	ModeFiles   string = "files"
	ModeTypes   string = "types"
	ModeSample  string = "sample"
	MemberLimit int    = 10000
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
)
