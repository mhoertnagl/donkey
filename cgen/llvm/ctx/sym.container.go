package ctx

import "github.com/llir/llvm/ir/value"

type SymbolTable map[string]Symbol

type Scopes []SymbolTable

type SymbolContainer struct {
	parent *SymbolContainer
	scopes Scopes
}

func (c *SymbolContainer) InitSymbolContainer(parent *SymbolContainer) {
	c.parent = parent
	c.scopes = make(Scopes, 0)
	c.PushScope()
}

func (c *SymbolContainer) PushScope() {
	c.scopes = append(c.scopes, make(SymbolTable))
}

func (c *SymbolContainer) PopScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *SymbolContainer) Set(name string, value value.Value) {
	c.scopes[len(c.scopes)-1][name] = NewValueSymbol(value)
}

func (c *SymbolContainer) Get(name string) (Symbol, bool) {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if sym, ok := c.scopes[i][name]; ok {
			return sym, true
		}
	}
	if c.parent != nil {
		return c.parent.Get(name)
	}
	return nil, false
}
