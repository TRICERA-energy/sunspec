package sunspec

// Symbol defines an element in the enumeration of a point.
type Symbol interface {
	Name() string
	Value() uint32
}

// SymbolDef is the definition of a sunspec symbol element.
type SymbolDef struct {
	Meta
	Name  string `json:"name"`
	Value uint32 `json:"value"`
}

func (def *SymbolDef) Simplify() {
	def.Meta.Simplify()
}

type symbol struct {
	name  string
	value uint32
}

func (s *symbol) Name() string { return s.name }

func (s *symbol) Value() uint32 { return s.value }

type Symbols map[uint32]Symbol

// Symbol retrieves the first symbol from the collection, identified by the given name.
func (sym Symbols) Symbol(name string) Symbol {
	for _, s := range sym {
		if s.Name() == name {
			return s
		}
	}
	return nil
}

// Symbols retrieves all symbols from the collection, identified by the given name.
func (sym Symbols) Symbols(names ...string) Symbols {
	if len(names) == 0 {
		return sym
	}
	col := make(Symbols, len(names))
	for key, s := range sym {
		for _, name := range names {
			if s.Name() == name {
				col[key] = s
				break
			}
		}
	}
	return col
}
