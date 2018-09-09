package main

import (
	"os"

	"github.com/mhoertnagl/donkey/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
