package sunspec

import (
	"errors"
	"regexp"
)

// Model
type Model interface {
	// Group defines a sunspec container for points.
	Group
	// ID returns the models identifier as defined by the first point "ID".
	ID() Uint16
	// Length returns the model length as defined by the second point "L".
	Length() Uint16
}

// Definition describes the behavior of a model reference, which can be instantiated.
type Definition interface {
	// ID returns the model´s identifier.
	ID() uint16
	// Instance derives a new useable model from the definition.
	Instance(adr uint16, callback func(pts []Point) error) (Model, error)
}

// ModelDef is the definition of a sunspec Model.
type ModelDef struct {
	Id          uint16   `json:"id"`
	Group       GroupDef `json:"group"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"desc,omitempty"`
	Detail      string   `json:"detail,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Comments    []string `json:"comments,omitempty"`
}

var _ Definition = (*ModelDef)(nil)

// ID returns the definitions model identifier.
func (def *ModelDef) ID() uint16 {
	return def.Id
}

// Instance derives a new useable Model from the defintion.
func (def *ModelDef) Instance(adr uint16, callback func(pts []Point) error) (Model, error) {
	m := &model{}

	var iterate func(def GroupDef) (Group, error)

	iterate = func(def GroupDef) (Group, error) {
		g := &group{
			name:   def.Name,
			atomic: bool(def.Atomic),
		}
		if m.group == nil {
			m.group = g
		}
		for _, def := range def.Points {
			for c := m.count(def.Count); c != 0; c-- {
				var p Point
				b := point{
					name:     def.Name,
					static:   bool(def.Static),
					writable: bool(def.Writable),
					address:  adr,
				}
				f := scale{def.ScaleFactor}
				s := make(Symbols, len(def.Symbols))
				for _, sym := range def.Symbols {
					s[sym.Value] = &symbol{sym.Name, sym.Value}
				}
				switch def.Type {
				case "int16":
					p = &t_Int16{b, toInt16(def.Value), f}
				case "int32":
					p = &t_Int32{b, toInt32(def.Value), f}
				case "int64":
					p = &t_Int64{b, toInt64(def.Value), f}
				case "pad":
					p = &t_Pad{b}
				case "sunnsf":
					p = &t_Sunssf{b, toInt16(def.Value)}
				case "uint16":
					p = &t_Uint16{b, toUint16(def.Value), f}
				case "uint32":
					p = &t_Uint32{b, toUint32(def.Value), f}
				case "uint64":
					p = &t_Uint64{b, toUint64(def.Value), f}
				case "acc16":
					p = &t_Acc16{b, toUint16(def.Value), f}
				case "acc32":
					p = &t_Acc32{b, toUint32(def.Value), f}
				case "acc64":
					p = &t_Acc64{b, toUint64(def.Value), f}
				case "bitfield16":
					p = &t_Bitfield16{b, toUint16(def.Value), s}
				case "bitfield32":
					p = &t_Bitfield32{b, toUint32(def.Value), s}
				case "bitfield64":
					p = &t_Bitfield64{b, toUint64(def.Value), s}
				case "enum16":
					p = &t_Enum16{b, toUint16(def.Value), s}
				case "enum32":
					p = &t_Enum32{b, toUint32(def.Value), s}
				case "string":
					p = &t_String{b, append(make([]byte, 0, def.Size*2), toByteS(def.Value)...)}
				case "float32":
					p = &t_Float32{b, toFloat32(def.Value)}
				case "float64":
					p = &t_Float64{b, toFloat64(def.Value)}
				case "ipaddr":
					p = &t_Ipaddr{b, [4]byte{}} // initial value ToDo
				case "ipv6addr":
					p = &t_Ipv6addr{b, [16]byte{}} // initial value ToDo
				case "eui48":
					p = &t_Eui48{b, [8]byte{}} // initial value ToDo
				}
				g.points = append(g.points, p)
				adr += p.Quantity()
			}
		}
		// ToDo: refactor initialization of points
		for _, p := range g.points {
			if i, ok := p.(interface{ init(p Point) }); ok {
				i.init(p)
			}
		}
		if callback != nil {
			if err := callback(g.points); err != nil {
				return nil, err
			}
		}
		for _, def := range def.Groups {
			for c := m.count(def.Count); c != 0; c-- {
				x, err := iterate(def)
				if err != nil {
					return nil, err
				}
				g.groups = append(g.groups, x)
			}
		}
		return g, nil
	}

	if _, err := iterate(def.Group); err != nil {
		return nil, err
	}

	m.ID().Set(def.Id)
	m.Length().Set(m.Quantity() - 2)

	return m, nil
}

// model is internally used to build out a usable model.
type model struct{ *group }

// count returns the number of occurrences of a point or group in the model.
func (m *model) count(c interface{}) uint16 {
	switch v := c.(type) {
	case int:
		return uint16(v)
	case float64:
		return uint16(v)
	case string:
		for _, p := range m.Points() {
			if p.Name() == v {
				switch p := p.(type) {
				case Int16:
					return uint16(p.Get())
				case Int32:
					return uint16(p.Get())
				case Int64:
					return uint16(p.Get())
				case Uint16:
					return uint16(p.Get())
				case Uint32:
					return uint16(p.Get())
				case Uint64:
					return uint16(p.Get())
				case Acc16:
					return uint16(p.Get())
				case Acc32:
					return uint16(p.Get())
				case Acc64:
					return uint16(p.Get())
				}
			}
		}
	}
	return 1
}

// ID returns the models identifier as defined by the first point "ID".
func (m *model) ID() Uint16 {
	if id := m.Points().Point("ID"); id != nil {
		return id.(Uint16)
	}
	return nil
}

// Length returns the model length as defined by the second point "L".
func (m *model) Length() Uint16 {
	if l := m.Points().Point("L"); l != nil {
		return l.(Uint16)
	}
	return nil
}

// Verify validates the given model, checking for its compliance regarding the official sunspec specification.
func Verify(m Model) error {
	if m.Length().Get()+2 != m.Quantity() {
		return errors.New("sunspec: Identifier L does not correlate with model quantity")
	}
	adr := m.Address()
	// spec ref 4.2.1 "An ID MUST consist of only alphanumeric characters
	// and the underscore character" - applies to group, point and symbol
	r, _ := regexp.Compile("^([[:alnum:]]|_)+$")
	return iterate(m, func(g Group) error {
		switch {
		case g.Address() != adr:
			return errors.New("sunspec: the given address range is not continuos")
		case !r.Match([]byte(g.Name())):
			return errors.New("sunspec: the name is violating the specifications definition")
		case g.Points() == nil:
			return errors.New("sunspec: the group is missing it´s point definition")
		}
		for _, p := range g.Points() {
			switch {
			case p.Address() != adr:
				return errors.New("sunspec: the given address range is not continuos")
			case !r.Match([]byte(p.Name())):
				return errors.New("sunspec: the name is violating the specifications definition")
			}
			adr += p.Quantity()
		}
		return nil
	})
}

// Models is a collection wrapper for multiple models.
// Offering functionalities applicable for them.
type Models []Model

// First returns the first model from the collection.
func (mls Models) First() Model { return mls[0] }

// Last returns the last model from the collection.
func (mls Models) Last() Model { return mls[len(mls)-1] }

// Model returns the first model found from the collection identified by the given id.
// If no model is found nil is returned instead.
func (mls Models) Model(id uint16) Model {
	for _, m := range mls {
		if m.ID().Get() == id {
			return m
		}
	}
	return nil
}

// Model returns a sub-collection of models identified by the given ids from the container.
// If ids is omitted a copy of the container itself is returned.
func (mls Models) Models(ids ...uint16) Models {
	if len(ids) == 0 {
		return append(Models(nil), mls...)
	}
	col := make(Models, 0, len(ids))
	for _, m := range mls {
		for _, id := range ids {
			if m.ID().Get() == id {
				col = append(col, m)
				break
			}
		}
	}
	return col
}

// Index returns the merged indexes of all models in the collection.
func (mls Models) Index() []Index {
	idx := make([]Index, 0, len(mls))
	for _, m := range mls {
		idx = append(idx, m)
	}
	return merge(idx)
}
