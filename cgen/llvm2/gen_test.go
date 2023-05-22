package llvm2_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/mhoertnagl/donkey/cgen/llvm2"
	"github.com/mhoertnagl/donkey/lexer"
	"github.com/mhoertnagl/donkey/parser"
	"github.com/mhoertnagl/donkey/utils/fs"
)

// TODO: assignments
// TODO: tests assignments
// TODO: tests nested if

func TestCodeGeneration(t *testing.T) {
	files := fs.FindFilesWithExtension("tests", ".dk")
	for _, actFile := range files {
		println(actFile)
		expFile := strings.ReplaceAll(actFile, ".dk", ".ll")
		act := compile(t, actFile)
		expBin, _ := os.ReadFile(expFile)
		exp := string(expBin)
		if !equalsIgnoreSpace(exp, act) {
			t.Errorf("Expected [%s] but got [%s]", exp, act)
		}
	}
}

func compile(t *testing.T, file string) string {
	t.Helper()
	input, _ := os.ReadFile(file)
	lexer := lexer.NewLexer(string(input))
	parser := parser.NewParser(lexer)
	prog := parser.Parse()
	gen := llvm2.NewLlvmCodegen()
	return gen.Generate(prog)
}

func equalsIgnoreSpace(exp, act string) bool {
	matcher := regexp.MustCompile(`\s+`)
	exp = matcher.ReplaceAllString(exp, "")
	act = matcher.ReplaceAllString(act, "")
	return exp == act
}
