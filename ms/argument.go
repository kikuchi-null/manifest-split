package ms

import (
	"flag"
	"log"
)

// mode: default, types, files
type Args struct {
	Input  string
	Output string
	Mode   string
	Num    int
}

// func RecieveArgs() (a Args) {

// 	fmt.Print("Path to the input package.xml: ")
// 	fmt.Scan(&a.Input)

// 	fmt.Print("Path to the output directory: ")
// 	fmt.Scan(&a.Output)

// 	fmt.Print("Split Mode(Default, ): ")
// 	fmt.Scan(&a.Mode)

// 	fmt.Print("Number of files to split into: ")
// 	fmt.Scan(&a.Num)

// 	a.validate()

// 	fmt.Println(a)

// 	return

// }

func RecieveArgs() (a Args) {
	i := flag.String("input", "", "分割したいpackage.xmlのパス")
	o := flag.String("output", "", "出力先のパス")
	m := flag.String("mode", "default", "分割モード（任意）")
	n := flag.Int("n", 1, "1ファイルに含まれるコンポーネント数の上限(最大1万) または 分割したいファイル数")
	flag.Parse()

	a = Args{
		Input:  *i,
		Output: *o,
		Mode:   *m,
		Num:    *n,
	}

	// 入力値の検証
	a.validate()

	return
}

func (a *Args) validate() (err error) {
	if a.Mode != ModeSample && (a.Input == "" || a.Output == "") {
		flag.Usage()
		log.Fatalf("inputとoutputを指定してください")
	}

	if a.Mode == ModeDefault && a.Num > MemberLimit {
		flag.Usage()
		log.Fatalf("コンポーネント数が10000以上だと組織からメタデータを取得できません")
	}

	if a.Mode == ModeDefault && a.Num < 1 {
		flag.Usage()
		log.Fatalf("ファイル数は1以上を指定してください")
	}

	return
}
