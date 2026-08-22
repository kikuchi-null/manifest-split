package ms

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
)

// 入力を格納する
type Args struct {
	Input  string
	Output string
	Mode   string
	Num    int
	Clean  bool
}

// コマンドラインフラグ
// flag.Parse()自体はmain.goで一度だけ呼び出す
var (
	fMode   = flag.String("mode", "", "分割モード (default, files, types, sample)")
	fInput  = flag.String("input", "", "分割したいpackage.xmlのパス")
	fOutput = flag.String("output", "", "出力先ディレクトリ")
	fNum    = flag.Int("num", 0, "1ファイルに含まれるコンポーネント数 または 分割したいファイル数")
	fClean  = flag.Bool("clean", false, "出力先ディレクトリの既存のpackage.xmlを削除してから書き込む")
)

// mode/input/output/numのいずれかがフラグで指定されていれば非対話(バッチ)モードとする
func isBatchMode() bool {
	return *fMode != "" || *fInput != "" || *fOutput != "" || *fNum != 0
}

// 入力を受け取る
// フラグ指定があれば非対話で、なければ従来通りターミナルから対話的に受け取る
func RecieveArgs() (a Args) {

	if isBatchMode() {
		a.Mode = *fMode
		a.Input = *fInput
		a.Output = *fOutput
		a.Num = *fNum
		a.Clean = *fClean

		a.validate()

		return
	}

	fmt.Println("====== manifest-split ======")
	fmt.Println("入力方法: 入力したらEnter")

	fmt.Print("モード選択(default, files, types, sample): ")
	fmt.Scanln(&a.Mode)

	if a.Mode != ModeSample {
		// サンプル作成モード以外
		fmt.Print("分割したいpackage.xml: ")
		fmt.Scanln(&a.Input)
	}

	fmt.Print("出力先: ")
	fmt.Scanln(&a.Output)

	if a.Mode == ModeDefault || a.Mode == ModeFiles || a.Mode == "" {
		fmt.Print("1ファイルに含まれるコンポーネント数(1〜10000) または 分割したいファイル数: ")
		fmt.Scanln(&a.Num)
	}

	// -cleanは対話モードでも独立して有効
	a.Clean = *fClean

	a.validate()

	return

}

// 入力検証
func (a *Args) validate() {

	// 入力不備があった場合はtrue
	var isError bool
	var errMessages []string

	// 入力検証
	if a.Mode == "" {
		// モード指定がない場合はデフォルトモードで起動
		a.Mode = ModeDefault
	}

	if result := slices.Contains(GetModes(), a.Mode); !result {
		errMessages = append(errMessages, "Error: モードは[default, files, types, sample]から選択してください")
		isError = true
	}

	if a.Mode != ModeSample && (a.Input == "" || a.Output == "") {
		// 入力ファイルと出力先の入力確認
		errMessages = append(errMessages, "Error: 分割対象と出力先を指定してください")
		isError = true
	}

	if a.Mode == ModeSample && a.Output == "" {
		// 出力先の入力確認
		errMessages = append(errMessages, "Error: 出力先を指定してください")
		isError = true
	}

	if a.Mode == ModeDefault && (a.Num < 1 || a.Num > MemberLimit) {
		// xmlファイルに含まれるコンポーネント数上限・下限確認
		errMessages = append(errMessages, "Error: コンポーネント数は1〜10000までで指定してください")
		isError = true
	}

	if a.Mode == ModeFiles && a.Num < 1 {
		// ファイル数
		errMessages = append(errMessages, "Error: ファイル数は1以上を指定してください")
		isError = true
	}

	if isError {
		color.Red("\n====== エラー ======")
		for _, mss := range errMessages {
			color.Red(mss)
		}
		os.Exit(1)
	}

}
