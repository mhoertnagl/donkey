package ctx

type Symbols map[string]Symbol

type Context interface {
	Get(name string) Symbol
}
