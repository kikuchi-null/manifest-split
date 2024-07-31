package ms

import (
	"flag"
	"fmt"
	"log"
)

// 入力を格納する
type Args struct {
	Input  string
	Output string
	Mode   string
	Num    int
}

// ターミナルからの入力を受け取る
func RecieveArgs() (a Args) {

	fmt.Println("入力したらEnter")
	fmt.Print("分割したいpackage.xmlのパス: ")
	fmt.Scanln(&a.Input)

	fmt.Print("出力先のパス: ")
	fmt.Scanln(&a.Output)

	fmt.Print("分割モード(Enterでデフォルトモード): ")
	fmt.Scanln(&a.Mode)

	fmt.Print("1ファイルに含まれるコンポーネント数の上限(最大1万) または 分割したいファイル数: ")
	fmt.Scanln(&a.Num)

	a.validate()

	return

}

// func RecieveArgs() (a Args) {
// 	i := flag.String("input", "", "分割したいpackage.xmlのパス")
// 	o := flag.String("output", "", "出力先のパス")
// 	m := flag.String("mode", "default", "分割モード（任意）")
// 	n := flag.Int("n", 1, "1ファイルに含まれるコンポーネント数の上限(最大1万) または 分割したいファイル数")
// 	flag.Parse()

// 	a = Args{
// 		Input:  *i,
// 		Output: *o,
// 		Mode:   *m,
// 		Num:    *n,
// 	}

// 	// 入力値の検証
// 	a.validate()

// 	return
// }

func (a *Args) validate() (err error) {

	// 入力検証
	if a.Mode == "" {
		// モード指定がない場合はデフォルトモードで起動
		a.Mode = ModeDefault
	}

	if a.Mode == ModeSample && a.Output == "" {
		// 出力先の入力確認
		flag.Usage()
		log.Fatalf("outputを指定してください")
	}

	if a.Mode != ModeSample && (a.Input == "" || a.Output == "") {
		// 入力ファイルと出力先の入力確認
		flag.Usage()
		log.Fatalf("inputとoutputを指定してください")
	}

	if a.Mode == ModeDefault && (a.Num < 1 || a.Num > MemberLimit) {
		// xmlファイルに含まれるコンポーネント数上限・下限確認
		flag.Usage()
		log.Fatalf("コンポーネント数は1〜10000までで指定してください")
	}

	if a.Mode == ModeFiles && a.Num < 1 {
		// 1ファイルに含まれる
		flag.Usage()
		log.Fatalf("ファイル数は1以上を指定してください")
	}

	return
}
