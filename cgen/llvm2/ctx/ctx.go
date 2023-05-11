package ctx

type Symbols map[string]Symbol

type Scopes []Symbols

type Context interface {
	Get(name string) Symbol
}
