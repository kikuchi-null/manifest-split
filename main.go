package main

import (
	"fmt"
	"manifest-split/ms"
)

func main() {

	var input string
	fmt.Print("input: ")
	fmt.Scan(&input)
	fmt.Println(input)

	args := ms.RecieveArgs()

	manifest := ms.ReadXML(args.Input)

	if args.Mode == ms.ModeFiles {
		manifest.GenerateModeFileSize(args.Output, args.Num)
	} else if args.Mode == ms.ModeTypes {
		manifest.GenerateModeTypes(args.Output)
	}

}
