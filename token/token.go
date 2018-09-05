package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	ID      = "ID"
	INT     = "INT"
	ASSIGN  = "="
	PLUS    = "+"
	COMMA   = ","
	SCOLON  = ";"
	LPAR    = "("
	RPAR    = ")"
	LBRA    = "{"
	RBRA    = "}"
	FUN     = "FUN"
	LET     = "LET"
)
