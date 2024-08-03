package main

import (
	"manifest-split/ms"
)

func main() {
	run()
}

func run() {

	// ターミナルからの入力を受け取る
	args := ms.RecieveArgs()

	// 出力先 ディレクトリの生成
	ms.GenerateOutputDirectory(args.Output)

	switch args.Mode {
	case ms.ModeSample:
		ms.GenerateLargePackageXML(args.Output)
	case ms.ModeTypes:
		// Typesごとに分割
		manifest := ms.ReadXML(args.Input)
		manifest.GenerateXMLModeTypes(args.Output)
	default:
		// defalt または files
		// package.xmlに含まれるコンポーネント数が10000以下になるように分割
		manifest := ms.ReadXML(args.Input)
		manifest.GenerateXML(args.Output, args.Mode, args.Num)
	}
}
