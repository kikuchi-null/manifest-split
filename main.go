package main

import (
	"manifest-split/msplit"
)

func main() {

	args := msplit.RecieveArgs()

	manifest := msplit.ReadXML(args.Input)

	if args.Mode == "files" {
		manifest.GenerateModeFileSize(args.Output, args.Num)
	} else if args.Mode == "types" {
		manifest.GenerateModeTypes(args.Output)
	}

}
