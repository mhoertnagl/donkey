package parser

import (
	"bytes"
	"fmt"
	"strings"
)

const INDENT = 3

func PrintParseTree(node Node) string {
	var buf bytes.Buffer
	return printParseTree(0, &buf, node)
}

func printParseTree(indent int, buf *bytes.Buffer, node Node) string {
	switch n := node.(type) {
	case *Program:
		for _, stmt := range n.Statements {
			printParseTree(indent, buf, stmt)
			buf.WriteString("\n")
		}
	case *BlockStatement:
		len := len(n.Statements) - 1
		buf.WriteString("BLOCK\n")
		for index := 0; index < len; index++ {
			printIntermediate(indent, buf, n.Statements[index])
		}
		printFinal(indent, buf, n.Statements[len])
	case *LetStatement:
		buf.WriteString(fmt.Sprintf("LET %s\n", n.Name))
		printFinal(indent, buf, n.Value)
	case *ReturnStatement:
		buf.WriteString("RETURN\n")
		printFinal(indent, buf, n.Value)
	case *FunDefStatement:
		buf.WriteString(fmt.Sprintf("FUN%s\n", n.Params))
		printFinal(indent, buf, n.Body)
	// case *FunctionLiteral:
	//   buf.WriteString(fmt.Sprintf("FUN%s\n", n.Params))
	//   printFinal(indent, buf, n.Body)
	case *BinaryExpression:
		buf.WriteString(fmt.Sprintf("INFIX(%s)\n", n.Operator))
		printIntermediate(indent, buf, n.Left)
		printFinal(indent, buf, n.Right)
	case *PrefixExpression:
		buf.WriteString(fmt.Sprintf("PREFIX(%s)\n", n.Operator))
		printFinal(indent, buf, n.Value)
	case *Identifier:
		buf.WriteString(n.String())
	case *Integer:
		buf.WriteString(n.String())
	case *Boolean:
		buf.WriteString(n.String())
	}
	return buf.String()
}

func printIntermediate(indent int, buf *bytes.Buffer, n Node) {
	buf.WriteString(strings.Repeat(" ", indent*3))
	buf.WriteString(" ├ ")
	printParseTree(indent+1, buf, n)
	buf.WriteString("\n")
}

func printFinal(indent int, buf *bytes.Buffer, n Node) {
	buf.WriteString(strings.Repeat(" ", indent*3))
	buf.WriteString(" └ ")
	printParseTree(indent+1, buf, n)
	//buf.WriteString("\n")
}
