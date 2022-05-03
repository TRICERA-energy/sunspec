package sunspec

// Point defines the generic behavior all sunspec types have in common.
type Point interface {
	// Index defines the locality of the point in a modbus address space.
	Index
	// Name returns the point´s identifier.
	Name() string
	// Valid specifies whether the underlying value is implemented by the device.
	Valid() bool
	// Origin returns the point´s associated group
	Origin() Group
	// Static specifies whether the point is expected to stay constant - not change over time.
	Static() bool
	// Writable specifies whether the point can be written to.
	Writable() bool
	// encode puts the point´s value into a buffer.
	encode(buf []byte) error
	// decode sets the point´s value from a buffer.
	decode(buf []byte) error
}

// PointDef is the definition of a sunspec point element.
type PointDef struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value,omitempty"`
	Count       interface{} `json:"count,omitempty"`
	Size        uint16      `json:"size"`
	ScaleFactor interface{} `json:"sf,omitempty"`
	Units       string      `json:"units,omitempty"`
	Writable    writable    `json:"access,omitempty"`
	Mandatory   mandatory   `json:"mandatory,omitempty"`
	Static      static      `json:"static,omitempty"`
	Label       string      `json:"label,omitempty"`
	Description string      `json:"desc,omitempty"`
	Detail      string      `json:"detail,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	Comments    []string    `json:"comments,omitempty"`
	Symbols     []SymbolDef `json:"symbols,omitempty"`
}

func (def *PointDef) Instance(adr uint16, o Group) Point {
	p := point{
		name:     def.Name,
		static:   bool(def.Static),
		writable: bool(def.Writable),
		origin:   o,
		address:  adr,
	}
	f := scale{def.ScaleFactor}
	s := make(Symbols, len(def.Symbols))
	for _, sym := range def.Symbols {
		s[sym.Value] = &symbol{sym.Name, sym.Value}
	}

	init := map[string]func() Point{
		"int16":      func() Point { return &tInt16{p, toInt16(def.Value), f} },
		"int32":      func() Point { return &tInt32{p, toInt32(def.Value), f} },
		"int64":      func() Point { return &tInt64{p, toInt64(def.Value), f} },
		"pad":        func() Point { return &tPad{p} },
		"sunssf":     func() Point { return &tSunssf{p, toInt16(def.Value)} },
		"uint16":     func() Point { return &tUint16{p, toUint16(def.Value), f} },
		"uint32":     func() Point { return &tUint32{p, toUint32(def.Value), f} },
		"unit64":     func() Point { return &tUint64{p, toUint64(def.Value), f} },
		"acc16":      func() Point { return &tAcc16{p, toUint16(def.Value), f} },
		"acc32":      func() Point { return &tAcc32{p, toUint32(def.Value), f} },
		"acc64":      func() Point { return &tAcc64{p, toUint64(def.Value), f} },
		"count":      func() Point { return &tCount{p, toUint16(def.Value)} },
		"bitfield16": func() Point { return &tBitfield16{p, toUint16(def.Value), s} },
		"bitfield32": func() Point { return &tBitfield32{p, toUint32(def.Value), s} },
		"bitfield64": func() Point { return &tBitfield64{p, toUint64(def.Value), s} },
		"enum16":     func() Point { return &tEnum16{p, toUint16(def.Value), s} },
		"enum32":     func() Point { return &tEnum32{p, toUint32(def.Value), s} },
		"string":     func() Point { return &tString{p, append(make([]byte, 0, def.Size*2), toByteS(def.Value)...)} },
		"float32":    func() Point { return &tFloat32{p, toFloat32(def.Value)} },
		"float64":    func() Point { return &tFloat64{p, toFloat64(def.Value)} },
		"ipaddr":     func() Point { return &tIpaddr{p, [4]byte{}} },    // initial value ToDo
		"ipv6addr":   func() Point { return &tIpv6addr{p, [16]byte{}} }, // initial value ToDo
		"eui48":      func() Point { return &tEui48{p, [8]byte{}} },     // initial value ToDo
	}
	return init[def.Type]()
}

// point is internally used to build out a useable model
type point struct {
	name     string
	origin   Group
	static   bool
	writable bool
	address  uint16
}

// Address returns the modbus starting address of the point.
func (p *point) Address() uint16 { return p.address }

// ID returns the point´s identifier
func (p *point) Name() string { return p.name }

// Writable specifies whether the point can be written to.
func (p *point) Writable() bool { return p.writable }

// Origin returns the point´s associated group
func (p *point) Origin() Group { return p.origin }

// Static specifies whether the points underlying data is supposed to be constant,
// meaning it is not supposed to change over time.
func (p *point) Static() bool { return p.static }

// Points is a collection wrapper for multiple Points.
// Offering functionalities applicable for them.
type Points []Point

// First returns the first point of the collection
func (pts Points) First() Point { return pts[0] }

// Last returns the last point of the collection
func (pts Points) Last() Point { return pts[len(pts)-1] }

// Quantity returns the total number of registers (2-Byte-Tuples/words)
// required to store the point in a modbus address space.
func (pts Points) Quantity() uint16 {
	var l uint16
	for _, p := range pts {
		l += p.Quantity()
	}
	return l
}

// Point returns the first immediate point identified by name.
func (pts Points) Point(name string) Point {
	for _, p := range pts {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// Points returns all immediate points identified by names.
// If names are omitted all points are returned.
func (pts Points) Points(names ...string) Points {
	if len(names) == 0 {
		return append(Points(nil), pts...)
	}
	col := make(Points, 0, len(names))
	for _, p := range pts {
		for _, id := range names {
			if p.Name() == id {
				col = append(col, p)
				break
			}
		}
	}
	return col
}

// address is internally used to get the address of a continuous collection of points.
func (pts Points) address() uint16 { return pts[0].Address() }

// Index returns the merged indexes of all points in the collection.
func (pts Points) Index() []Index {
	idx := make([]Index, 0, len(pts))
	for _, p := range pts {
		idx = append(idx, p)
	}
	return merge(idx)
}

// index is internally used to get the locality of a continuous collection of points.
func (pts Points) index() Index {
	return index{address: pts.address(), quantity: pts.Quantity()}
}

// decode sets the value for all points in the collection as stored in the buffer.
func (pts Points) decode(buf []byte) error {
	for _, p := range pts {
		if err := p.decode(buf); err != nil {
			return err
		}
		buf = buf[2*p.Quantity():]
	}
	return nil
}

// encode puts the values of the points in the collection into the buffer.
func (pts Points) encode(buf []byte) error {
	for _, p := range pts {
		if err := p.encode(buf); err != nil {
			return err
		}
		buf = buf[2*p.Quantity():]
	}
	return nil
}
