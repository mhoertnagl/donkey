package aegis

import (
  "io"
  "fmt"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors

const fgEscSeq = "\x1b[38;2;%d;%d;%dm";

// type cfmt struct {
//   fgc *color
// }
// 
// func NewCfmt() *cfmt {
//   return &cfmt{fgc: White()};
// }

func SetTextColor(c *color) {
  fmt.Printf(fgEscSeq, c.r, c.g, c.b);
}

func FsetTextColor(w io.Writer, c *color) {
  fmt.Fprintf(w, fgEscSeq, c.r, c.g, c.b);
}
