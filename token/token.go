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
	COMMA   = ","
	SCOLON  = ";"
	LPAR    = "("
	RPAR    = ")"
	LBRA    = "{"
	RBRA    = "}"
	FUN     = "FUN"
	LET     = "LET"
)

var keywords = map[string]TokenType{
	"fun": FUN,
	"let": LET,
}

func LookupId(id string) TokenType {
	if tok, ok := keywords[id]; ok {
		return tok
	}
	return ID
}
