package parser

import (
	"testing"
)

func TestPrintParseTreeIdentifier(t *testing.T) {
  n := &Identifier{Value: "x"}
  expected := `x`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeInteger(t *testing.T) {
  n := &Integer{Value: 42}
  expected := `42`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeBoolean(t *testing.T) {
  n := &Boolean{Value: true}
  expected := `true`
  testParseTree(t, n, expected)
}

func TestPrintParseTreePrefix(t *testing.T) {
  val := &Integer{Value: 42}
  n := &PrefixExpression{Operator: "-", Value: val}
  expected := `PREFIX(-)
 └ 42`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeInfix(t *testing.T) {
  left := &Integer{Value: 42}
  right := &Integer{Value: 43}
  n := &BinaryExpression{Operator: "+", Left: left, Right: right}
  expected := `INFIX(+)
 ├ 42
 └ 43`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeFunctionLiteral(t *testing.T) {
  params := []*Identifier{}
  params = append(params, &Identifier{Value: "a"})
  params = append(params, &Identifier{Value: "b"})
  
  stmts := []Statement{}
  val := &Integer{Value: 42}
  stmt := &ReturnStatement{Value: val}
  stmts = append(stmts, stmt)
  block := &BlockStatement{Statements: stmts}
  n := &FunctionLiteral{Params: params, Body: block}
  expected := `FUN[a b]
 └ BLOCK
    └ RETURN
       └ 42`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeReturn(t *testing.T) {
  val := &Integer{Value: 42}
  n := &ReturnStatement{Value: val}
  expected := `RETURN
 └ 42`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeLet(t *testing.T) {
  name := &Identifier{Value: "x"}
  val := &Integer{Value: 42}
  n := &LetStatement{Name: name, Value: val}
  expected := `LET x
 └ 42`
  testParseTree(t, n, expected)
}

func TestPrintParseTreeBlock(t *testing.T) {
  stmts := []Statement{}
  
  name1 := &Identifier{Value: "x"}
  val1 := &Integer{Value: 42}
  stmt1 := &LetStatement{Name: name1, Value: val1}
  stmts = append(stmts, stmt1)
  
  name2 := &Identifier{Value: "y"}
  val2 := &Integer{Value: 43}
  stmt2 := &LetStatement{Name: name2, Value: val2}
  stmts = append(stmts, stmt2)  
  
  val3 := &BinaryExpression{Operator: "+", Left: name1, Right: name2}
  stmt3 := &ReturnStatement{Value: val3}
  stmts = append(stmts, stmt3)
  
  n := &BlockStatement{Statements: stmts}
  expected := `BLOCK
 ├ LET x
    └ 42
 ├ LET y
    └ 43
 └ RETURN
    └ INFIX(+)
       ├ x
       └ y`
//   expected := `BLOCK
//  ├ LET x
//  │  └ 42
//  ├ LET y
//  │  └ 43
//  └ RETURN
//     └ INFIX(+)
//        ├ x
//        └ y
// `
  testParseTree(t, n, expected)
}

func testParseTree(t *testing.T, n Node, expected string) {
  actual := PrintParseTree(n)
  if actual != expected {
    t.Errorf("Expected [\n%s\n] but got [\n%s\n].", expected, actual)
  }
}
