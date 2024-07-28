package ms

import (
	"flag"
	"fmt"
	"log"
)

// mode: default, types, files
type Args struct {
	Input  string
	Output string
	Mode   string
	Num    int
}

func RecieveArgs() (a Args) {
	i := flag.String("input", "", "分割したいpackage.xmlのパス")
	o := flag.String("output", "", "分割したpackage.xml出力先のパス")
	m := flag.String("mode", "default", "分割モード（任意）")
	n := flag.Int("n", 1, "1ファイルに含まれるコンポーネント数の上限(max 10000) または 分割したいファイル数")
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

func RecieveArgsInterraction() (a Args) {

	fmt.Print("分割したいpackage.xmlのパス: ")
	fmt.Scan(&a.Input)

	fmt.Print("分割したpackage.xml出力先のパス: ")
	fmt.Scan(&a.Output)

	fmt.Print("分割モード（任意）: ")
	fmt.Scan(&a.Mode)

	fmt.Print("1ファイルに含まれるコンポーネント数の上限(最大1万) または 分割したいファイル数: ")
	fmt.Scan(&a.Num)

	// 入力値の検証
	a.validate()

	return

}

func (a *Args) validate() {
	if a.Mode != ModeSample && (a.Input == "" || a.Output == "") {
		flag.Usage()
		log.Fatalf("inputとoutputを指定してください")
	}

	if a.Mode == ModeDefault && a.Num > MemberLimit {
		log.Fatalf("コンポーネント数が10000以上だと組織からメタデータを取得できません")
	}

	if a.Mode == ModeDefault && a.Num < 1 {
		flag.Usage()
		log.Fatalf("ファイル数は1以上を指定してください")
	}

}
