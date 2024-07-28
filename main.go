package main

import (
	"manifest-split/ms"
)

func main() {

	args := ms.RecieveArgs()

	manifest := ms.ReadXML(args.Input)

	if args.Mode == ms.ModeFiles {
		manifest.GenerateModeFileSize(args.Output, args.Num)
	} else if args.Mode == ms.ModeTypes {
		manifest.GenerateModeTypes(args.Output)
	}

}
