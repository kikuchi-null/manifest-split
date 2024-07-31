package main

import (
	"manifest-split/ms"
)

func main() {
	run()
}

func run() (err error) {

	// ターミナルからの入力を受け取る
	args := ms.RecieveArgs()

	// 出力先 ディレクトリの生成
	ms.GenerateOutputDirectory(args.Output)

	// package.xml生成
	if args.Mode == ms.ModeSample {
		// サンプルデータ生成
		ms.GenerateLargePackageXML(args.Output)

	} else {
		// package.xmlの分割
		manifest := ms.ReadXML(args.Input)

		switch args.Mode {
		case ms.ModeFiles:
			// 指定されたファイル数に分割
			manifest.GenerateModeFileSize(args.Output, args.Num)
		case ms.ModeTypes:
			// Typesごとに分割
			manifest.GenerateModeTypes(args.Output)
		default:
			// デフォルトモード
			// package.xmlに含まれるコンポーネント数が10000以下になるように分割
			manifest.GenerateModeDefault(args.Output, args.Num)
		}
	}
	return

}
