package main

import (
	"manifest-split/msplit"
)

func main() {

	args := msplit.RecieveArgs()

	manifest := msplit.ReadXML(args.Input)

	if args.Mode == msplit.ModeFiles {
		manifest.GenerateModeFileSize(args.Output, args.Num)
	} else if args.Mode == msplit.ModeTypes {
		manifest.GenerateModeTypes(args.Output)
	}

}
