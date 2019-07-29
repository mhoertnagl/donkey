package token

type TokenType string

type Token struct {
	Typ     TokenType
	Literal string
  Line    int
  Col     int
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	ID      TokenType = "ID"
	INT     TokenType = "INT"
	ASSIGN  TokenType = "="
	PLUS    TokenType = "+"
	MINUS   TokenType = "-"
	TIMES   TokenType = "*"
	DIV     TokenType = "/"
	INV     TokenType = "~"
	AND     TokenType = "&"
	OR      TokenType = "|"
	XOR     TokenType = "^"
  // TODO: NOR
  // NOR TokenType = "~|"
	SLL     TokenType = "<<"
	SRL     TokenType = ">>"
	SRA     TokenType = ">>>"
	ROL     TokenType = "<<>"
	ROR     TokenType = "<>>"
	NOT     TokenType = "!"
	CONJ    TokenType = "&&"
	DISJ    TokenType = "||"
	EQU     TokenType = "=="
	NEQ     TokenType = "!="
	LT      TokenType = "<"
	LE      TokenType = "<="
	GT      TokenType = ">"
	GE      TokenType = ">="
	COMMA   TokenType = ","
	SCOLON  TokenType = ";"
	LPAR    TokenType = "("
	RPAR    TokenType = ")"
	LBRA    TokenType = "{"
	RBRA    TokenType = "}"
	FUN     TokenType = "FUN"
	LET     TokenType = "LET"
	TRUE    TokenType = "TRUE"
	FALSE   TokenType = "FALSE"
	IF      TokenType = "IF"
	ELSE    TokenType = "ELSE"
	RETURN  TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"fun":    FUN,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupId(id string) TokenType {
	if tok, ok := keywords[id]; ok {
		return tok
	}
	return ID
}
