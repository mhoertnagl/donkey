package main

import (
	"flag"
	"os"

	"github.com/mhoertnagl/donkey/console"
	"github.com/mhoertnagl/donkey/repl"
)

func main() {

	lexOnly := flag.Bool("l", false, "a bool")
	parseOnly := flag.Bool("p", false, "a bool")

	flag.Parse()

	cargs := console.Args{
		LexOnly:   *lexOnly,
		ParseOnly: *parseOnly,
	}

	if len(flag.Args()) == 0 {
		repl.Start(os.Stdin, os.Stdout, cargs)
	} else {

	}
}
