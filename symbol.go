package sunspec

type SymbolDef struct {
	Name        string   `json:"name"`
	Value       uint32   `json:"value"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"desc,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Comments    []string `json:"comments,omitempty"`
}

type Symbol interface {
	Name() string
	Value() uint32
}

type symbol struct {
	name  string
	value uint32
}

func (s *symbol) Name() string { return s.name }

func (s *symbol) Value() uint32 { return s.value }

type Symbols map[uint32]Symbol

func (sym Symbols) Symbol(name string) Symbol {
	for _, s := range sym {
		if s.Name() == name {
			return s
		}
	}
	return nil
}

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
