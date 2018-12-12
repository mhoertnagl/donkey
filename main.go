package main

import (
	"os"
  "flag"

  "github.com/mhoertnagl/donkey/console"
	"github.com/mhoertnagl/donkey/repl"
)

func main() {
  
  lexOnly := flag.Bool("l", false, "a bool")
  
  flag.Parse()
  
  cargs := console.Args{
    LexOnly: *lexOnly,
  }
  
  if len(flag.Args()) == 0 {
    repl.Start(os.Stdin, os.Stdout, cargs)
  } else {
    
  }	
}
