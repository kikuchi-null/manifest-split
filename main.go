package main

import (
	"log"
	"manifest-split/ms"
)

func main() {

	err := run()
	if err != nil {
		log.Fatalln(err)
	}

}

func run() (err error) {

	// ターミナルからの入力を受け取る
	args := ms.RecieveArgs()

	// 出力先 ディレクトリの生成
	err = ms.GenerateOutputDirectory(args.Output)
	if err != nil {
		return
	}

	switch args.Mode {
	case ms.ModeSample: // sample
		// 大量のコンポーネントを含むサンプルデータを作成する
		err = ms.GenerateLargePackageXML(args.Output)
		return

	case ms.ModeTypes: // types
		// Typesごとにpackage.xmlを分割する
		manifest, err := ms.ReadXML(args.Input)
		if err != nil {
			return err
		}

		err = manifest.GenerateXMLModeTypes(args.Output)
		return err

	default: // defalt または files
		// package.xmlに含まれるコンポーネント数が10000以下になるように分割
		manifest, err := ms.ReadXML(args.Input)
		if err != nil {
			return err
		}

		manifest.SplitTypes()
		err = manifest.GenerateXML(args.Output, args.Mode, args.Num)
		return err
	}

}
