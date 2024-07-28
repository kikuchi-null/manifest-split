package main

import (
	"fmt"
	"manifest-split/ms"
	"os"
)

func main() {
	err := run()
	if err != nil {
		os.Exit(2)
	}
}

func run() (err error) {

	// ターミナルからの入力を受け取る
	args, err := ms.RecieveArgs()
	if err != nil {
		return
	}

	fmt.Printf("mode: %v\n", args.Mode)

	// 出力先 ディレクトリの生成
	err = ms.GenerateOutputDirectory(args.Output)
	if err != nil {
		return
	}

	// package.xml生成
	if args.Mode == ms.ModeSample {
		// サンプルデータ生成
		ms.GenerateLargePackageXML(args.Output)

	} else {
		// package.xmlの分割
		manifest, _ := ms.ReadXML(args.Input)

		if args.Mode == ms.ModeFiles {
			// 指定されたファイル数に分割
			manifest.GenerateModeFileSize(args.Output, args.Num)

		} else if args.Mode == ms.ModeTypes {
			// Typesごとに分割
			manifest.GenerateModeTypes(args.Output)

		} else {
			// デフォルトモード
			// package.xmlに含まれるコンポーネント数が10000以下になるように分割
			manifest.GenerateModeDefault(args.Output, args.Num)

		}

	}
	return

}
