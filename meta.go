package sunspec

type Meta struct {
	Label       string   `json:"label,omitempty"`
	Description string   `json:"desc,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Comments    []string `json:"comments,omitempty"`
}

func (m *Meta) Simplify() {
	*m = Meta{}
}
