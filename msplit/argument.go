package msplit

import (
	"flag"
	"log"
)

// mode: default(10000コンポーネント以下に分割), types, files
type Args struct {
	Input  string
	Output string
	Mode   string
	Num    int
}

func RecieveArgs() (a Args) {
	i := flag.String("input", "", "Path to the input XML file")
	o := flag.String("output", "", "Path to the output directory")
	m := flag.String("mode", "", "The number of the output files")
	n := flag.Int("n", 1, "Number of files to split into")
	flag.Parse()

	a = Args{
		Input:  *i,
		Output: *o,
		Mode:   *m,
		Num:    *n,
	}

	a.validate()

	return
}

func (a *Args) validate() {
	if a.Input == "" || a.Output == "" {
		flag.Usage()
		log.Fatal("Both input and output parameters are required")
	}

	if a.Mode == "" {
		a.Mode = "default"
	}

	if a.Mode == ModeFiles && a.Num < 1 {
		flag.Usage()
		log.Fatal("The Number of files must be greater than 1")
	}
}
