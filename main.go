package main

import (
	"flag"
	"fmt"
	"log"
	"manifest-split/ms"
)

var showVersion = flag.Bool("version", false, "バージョンを表示して終了")

func main() {

	// ms パッケージで定義されたフラグも含めてここで一度だけパースする
	flag.Parse()

	if *showVersion {
		fmt.Println("manifest-split", ms.Version)
		return
	}

	if err := run(); err != nil {
		log.Fatalln(err)
	}

}

func run() error {

	// 入力を受け取る(フラグ指定があれば非対話、なければターミナルから対話的に受け取る)
	args := ms.RecieveArgs()

	// 出力先 ディレクトリの生成(非破壊的なのでクリーンアップより先に行ってよい)
	if err := ms.GenerateOutputDirectory(args.Output); err != nil {
		return fmt.Errorf("出力先ディレクトリ %q の作成に失敗しました: %w", args.Output, err)
	}

	var written []string

	switch args.Mode {
	case ms.ModeSample: // sample
		// 大量のコンポーネントを含むサンプルデータを作成する
		w, err := ms.GenerateLargePackageXML(args.Output)
		if err != nil {
			return fmt.Errorf("サンプルの生成に失敗しました: %w", err)
		}
		written = w

	case ms.ModeTypes: // types
		// Typesごとにpackage.xmlを分割する
		manifest, err := ms.ReadXML(args.Input)
		if err != nil {
			return fmt.Errorf("%q の読み込みに失敗しました: %w", args.Input, err)
		}

		w, err := manifest.GenerateXMLModeTypes(args.Output)
		if err != nil {
			return fmt.Errorf("package.xmlの生成に失敗しました: %w", err)
		}
		written = w

	default: // defalt または files
		// package.xmlに含まれるコンポーネント数が10000以下になるように分割
		manifest, err := ms.ReadXML(args.Input)
		if err != nil {
			return fmt.Errorf("%q の読み込みに失敗しました: %w", args.Input, err)
		}

		manifest.SplitTypes()
		w, err := manifest.GenerateXML(args.Output, args.Mode, args.Num)
		if err != nil {
			return fmt.Errorf("package.xmlの生成に失敗しました: %w", err)
		}
		written = w
	}

	if args.Clean {
		// 生成が全て成功した後に、今回書き込まなかった残存生成物だけを削除する
		if err := ms.CleanOutputDirectory(args.Output, written); err != nil {
			return fmt.Errorf("出力先ディレクトリ %q のクリーンアップに失敗しました: %w", args.Output, err)
		}
	}

	return nil

}
