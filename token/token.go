package token

type TokenType string

type Token struct {
	Typ     TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	ID      = "ID"
	INT     = "INT"
	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	TIMES   = "*"
	DIV     = "/"
	INV     = "~"
	AND     = "&"
	OR      = "|"
	XOR     = "^"
	SLL     = "<<"
	SRL     = ">>"
	SRA     = ">>>"
	ROL     = "<<>"
	ROR     = "<>>"
	NOT     = "!"
	CONJ    = "&&"
	DISJ    = "||"
	EQU     = "=="
	NEQ     = "!="
	LT      = "<"
	LE      = "<="
	GT      = ">"
	GE      = ">="
	COMMA   = ","
	SCOLON  = ";"
	LPAR    = "("
	RPAR    = ")"
	LBRA    = "{"
	RBRA    = "}"
	FUN     = "FUN"
	LET     = "LET"
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	IF      = "IF"
	ELSE    = "ELSE"
	RETURN  = "RETURN"
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
