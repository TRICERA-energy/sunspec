package sunspec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
)

// Scalable defines the behavior of a point type which may be scaled using the definition:
//	ScaledValue = PointValue * (10^ScaleFactor)
type Scalable interface {
	Scaled() bool
	Factor() int16
}

// scale is internally used to store a scale factor
type scale struct {
	f interface{}
}

// init initializes the scale by setting its constant value or
// backtracing the references point of type sunssf.
func (s *scale) init(p Point) {
	switch sf := s.f.(type) {
	case int:
		s.f = int16(sf)
	case float64:
		s.f = int16(sf)
	case string:
		for g := p.Origin(); g != nil; g = g.Origin() {
			for _, p := range g.Points() {
				if p.Name() == sf {
					if p, ok := p.(Sunssf); ok {
						s.f = p
					}
				}
			}
		}
	}
	s.f = int16(0)
}

// Scaled specifies whether the point is scaled using an optional factor.
func (s *scale) Scaled() bool {
	return s.f != nil
}

// Factor returns the scale value of the point.
func (s *scale) Factor() int16 {
	switch sf := s.f.(type) {
	case int16:
		return sf
	case Sunssf:
		return sf.Get()
	}
	return 1
}

// ****************************************************************************

// Int16 represents the sunspec type int16.
type Int16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int16) error
	// Get returns the point´s underlying value.
	Get() int16
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Int16 struct {
	point
	data int16
	scale
}

var _ Int16 = (*t_Int16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Int16) Valid() bool { return t.Get() != -0x8000 }

// String formats the point´s value as string.
func (t *t_Int16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Int16) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Int16) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, uint16(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Int16) Decode(buf []byte) error {
	return t.Set(int16(binary.BigEndian.Uint16(buf)))
}

// Set sets the point´s underlying value.
func (t *t_Int16) Set(v int16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Int16) Get() int16 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Int16) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Int32 represents the sunspec type int32.
type Int32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int32) error
	// Get returns the point´s underlying value.
	Get() int32
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Int32 struct {
	point
	data int32
	scale
}

var _ (Int32) = (*t_Int32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Int32) Valid() bool { return t.Get() != -0x80000000 }

// String formats the point´s value as string.
func (t *t_Int32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Int32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Int32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, uint32(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Int32) Decode(buf []byte) error {
	return t.Set(int32((binary.BigEndian.Uint32(buf))))
}

// Set sets the point´s underlying value.
func (t *t_Int32) Set(v int32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Int32) Get() int32 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Int32) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Int64 represents the sunspec type int64.
type Int64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v int64) error
	// Get returns the point´s underlying value.
	Get() int64
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Int64 struct {
	point
	data int64
	scale
}

var _ Int64 = (*t_Int64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Int64) Valid() bool { return t.Get() != -0x8000000000000000 }

// String formats the point´s value as string.
func (t *t_Int64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Int64) Quantity() uint16 { return 4 }

// Encode puts the point´s value into a buffer.
func (t *t_Int64) Encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, uint64(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Int64) Decode(buf []byte) error {
	return t.Set(int64(binary.BigEndian.Uint64(buf)))
}

// Set sets the point´s underlying value.
func (t *t_Int64) Set(v int64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Int64) Get() int64 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Int64) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Pad represents the sunspec type pad.
type Pad interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
}

type t_Pad struct {
	point
}

var _ Pad = (*t_Pad)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Pad) Valid() bool { return false }

// String formats the point´s value as string.
func (t *t_Pad) String() string { return "" }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Pad) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Pad) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, 0x8000)
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Pad) Decode(buf []byte) error { return nil }

// ****************************************************************************

// Sunssf represents the sunspec type sunssf.
type Sunssf interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Get returns the point´s underlying value.
	Get() int16
}

type t_Sunssf struct {
	point
	data int16
}

var _ Sunssf = (*t_Sunssf)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Sunssf) Valid() bool { return t.Get() != -0x8000 }

// String formats the point´s value as string.
func (t *t_Sunssf) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Sunssf) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Sunssf) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, uint16(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Sunssf) Decode(buf []byte) error {
	return t.set(int16(binary.BigEndian.Uint16(buf)))
}

// set sets the point´s underlying value.
func (t *t_Sunssf) set(v int16) error {
	if v < -10 || v > 10 {
		return errors.New("sunspec: value out of boundary")
	}
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Sunssf) Get() int16 { return t.data }

// ****************************************************************************

// Uint16 represents the sunspec type uint16.
type Uint16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Uint16 struct {
	point
	data uint16
	scale
}

var _ Uint16 = (*t_Uint16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Uint16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *t_Uint16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Uint16) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Uint16) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Uint16) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *t_Uint16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Uint16) Get() uint16 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Uint16) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Uint32 represents the sunspec type uint32.
type Uint32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Uint32 struct {
	point
	data uint32
	scale
}

var _ Uint32 = (*t_Uint32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Uint32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *t_Uint32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Uint32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Uint32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Uint32) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *t_Uint32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Uint32) Get() uint32 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Uint32) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Uint64 represents the sunspec type uint64.
type Uint64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
	// Value returns the scaled value as defined by the specification.
	Value() float64
}

type t_Uint64 struct {
	point
	data uint64
	scale
}

var _ Uint64 = (*t_Uint64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Uint64) Valid() bool { return t.Get() != 0xFFFFFFFFFFFFFFFF }

// String formats the point´s value as string.
func (t *t_Uint64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Uint64) Quantity() uint16 { return 4 }

// Encode puts the point´s value into a buffer.
func (t *t_Uint64) Encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Uint64) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *t_Uint64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Uint64) Get() uint64 { return t.data }

// Value returns the scaled value as defined by the specification.
func (t *t_Uint64) Value() float64 { return float64(t.Get()) * math.Pow10(int(t.Factor())) }

// ****************************************************************************

// Acc16 represents the sunspec type acc16.
type Acc16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
}

type t_Acc16 struct {
	point
	data uint16
	scale
}

var _ Acc16 = (*t_Acc16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Acc16) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *t_Acc16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Acc16) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Acc16) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Acc16) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *t_Acc16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Acc16) Get() uint16 { return t.data }

// ****************************************************************************

// Acc32 represents the sunspec type acc32.
type Acc32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
}

type t_Acc32 struct {
	point
	data uint32
	scale
}

var _ (Acc32) = (*t_Acc32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Acc32) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *t_Acc32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Acc32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Acc32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Acc32) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *t_Acc32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Acc32) Get() uint32 { return t.data }

// ****************************************************************************

// Acc64 represents the sunspec type acc64.
type Acc64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Scalable defines the behavior of a point type which may be scaled using the definition.
	Scalable
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
}

type t_Acc64 struct {
	point
	data uint64
	scale
}

var _ Acc64 = (*t_Acc64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Acc64) Valid() bool { return t.Get() != 0 }

// String formats the point´s value as string.
func (t *t_Acc64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Acc64) Quantity() uint16 { return 4 }

// Encode puts the point´s value into a buffer.
func (t *t_Acc64) Encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Acc64) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *t_Acc64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Acc64) Get() uint64 { return t.data }

// ****************************************************************************

// Bitfield16 represents the sunspec type bitfield16.
type Bitfield16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [16]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type t_Bitfield16 struct {
	point
	data    uint16
	symbols Symbols
}

var _ Bitfield16 = (*t_Bitfield16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Bitfield16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *t_Bitfield16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Bitfield16) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Bitfield16) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Bitfield16) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *t_Bitfield16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Bitfield16) Get() uint16 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *t_Bitfield16) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 15:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *t_Bitfield16) Field() (f [16]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *t_Bitfield16) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Bitfield32 represents the sunspec type bitfield32.
type Bitfield32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [32]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type t_Bitfield32 struct {
	point
	data    uint32
	symbols Symbols
}

var _ Bitfield32 = (*t_Bitfield32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Bitfield32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *t_Bitfield32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Bitfield32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Bitfield32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Bitfield32) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *t_Bitfield32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Bitfield32) Get() uint32 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *t_Bitfield32) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 31:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *t_Bitfield32) Field() (f [32]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *t_Bitfield32) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Bitfield64 represents the sunspec type bitfield64.
type Bitfield64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint64) error
	// Get returns the point´s underlying value.
	Get() uint64
	// Flip sets the bit at position pos, starting at 0, to the value of v.
	Flip(pos int, v bool) error
	// Field returns the individual bit values as bool array.
	Field() [64]bool
	// States returns all active enumerated states, correlating the bit value to its symbol.
	States() []string
}

type t_Bitfield64 struct {
	point
	data    uint64
	symbols Symbols
}

var _ Bitfield64 = (*t_Bitfield64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Bitfield64) Valid() bool { return t.Get() != 0xFFFFFFFFFFFFFFFF }

// String formats the point´s value as string.
func (t *t_Bitfield64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Bitfield64) Quantity() uint16 { return 4 }

// Encode puts the point´s value into a buffer.
func (t *t_Bitfield64) Encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Bitfield64) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint64(buf))
}

// Set sets the point´s underlying value.
func (t *t_Bitfield64) Set(v uint64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Bitfield64) Get() uint64 { return t.data }

// Flip sets the bit at position pos, starting at 0, to the value of v.
func (t *t_Bitfield64) Flip(pos int, v bool) error {
	switch {
	case pos < 0 || pos > 63:
		return errors.New("sunspec: out of bounds bit position ")
	case v:
		return t.Set(t.Get() | (1 << pos))
	}
	return t.Set(t.Get() &^ (1 << pos))
}

// Field returns the individual bit values as bool array.
func (t *t_Bitfield64) Field() (f [64]bool) {
	for v, b := t.Get(), 0; b < len(f); b++ {
		f[b] = v&(1<<b) != 0
	}
	return f
}

// States returns all active enumerated states, correlating the bit value to its symbol.
func (t *t_Bitfield64) States() (s []string) {
	if !t.Valid() {
		return nil
	}
	for i, v := range t.Field() {
		if v {
			s = append(s, t.symbols[uint32(i)].Name())
		}
	}
	return s
}

// ****************************************************************************

// Enum16 represents the sunspec type enum16.
type Enum16 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint16) error
	// Get returns the point´s underlying value.
	Get() uint16
	// State returns the currently active enumerated state.
	State() string
}

type t_Enum16 struct {
	point
	data    uint16
	symbols Symbols
}

var _ Enum16 = (*t_Enum16)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Enum16) Valid() bool { return t.Get() != 0xFFFF }

// String formats the point´s value as string.
func (t *t_Enum16) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Enum16) Quantity() uint16 { return 1 }

// Encode puts the point´s value into a buffer.
func (t *t_Enum16) Encode(buf []byte) error {
	binary.BigEndian.PutUint16(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Enum16) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint16(buf))
}

// Set sets the point´s underlying value.
func (t *t_Enum16) Set(v uint16) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Enum16) Get() uint16 { return t.data }

// State returns the currently active enumerated state.
func (t *t_Enum16) State() string { return t.symbols[uint32(t.Get())].Name() }

// ****************************************************************************

// Enum32 represents the sunspec type enum32.
type Enum32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v uint32) error
	// Get returns the point´s underlying value.
	Get() uint32
	// State returns the currently active enumerated state.
	State() string
}

type t_Enum32 struct {
	point
	data    uint32
	symbols Symbols
}

var _ Enum32 = (*t_Enum32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Enum32) Valid() bool { return t.Get() != 0xFFFFFFFF }

// String formats the point´s value as string.
func (t *t_Enum32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Enum32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Enum32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Enum32) Decode(buf []byte) error {
	return t.Set(binary.BigEndian.Uint32(buf))
}

// Set sets the point´s underlying value.
func (t *t_Enum32) Set(v uint32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Enum32) Get() uint32 { return t.data }

// State returns the currently active enumerated state.
func (t *t_Enum32) State() string { return t.symbols[t.Get()].Name() }

// ****************************************************************************

// String represents the sunspec type string.
type String interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v string) error
	// Get returns the point´s underlying value.
	Get() string
}

type t_String struct {
	point
	data []byte
}

var _ String = (*t_String)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_String) Valid() bool { return t.Get() != "" }

// String formats the point´s value as string.
func (t *t_String) String() string { return t.Get() }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_String) Quantity() uint16 { return uint16(cap(t.data) / 2) }

// Encode puts the point´s value into a buffer.
func (t *t_String) Encode(buf []byte) error {
	copy(buf, []byte(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_String) Decode(buf []byte) error {
	return t.Set(string(buf[:2*t.Quantity()]))
}

// Set sets the point´s underlying value.
func (t *t_String) Set(v string) error {
	copy(t.data[:cap(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *t_String) Get() string { return string(t.data) }

// ****************************************************************************

// Float32 represents the sunspec type float32.
type Float32 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v float32) error
	// Get returns the point´s underlying value.
	Get() float32
}

type t_Float32 struct {
	point
	data float32
}

var _ Float32 = (*t_Float32)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Float32) Valid() bool { return t.Get() != 0x7FC00000 }

// String formats the point´s value as string.
func (t *t_Float32) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Float32) Quantity() uint16 { return 2 }

// Encode puts the point´s value into a buffer.
func (t *t_Float32) Encode(buf []byte) error {
	binary.BigEndian.PutUint32(buf, math.Float32bits(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Float32) Decode(buf []byte) error {
	return t.Set(math.Float32frombits(binary.BigEndian.Uint32(buf)))
}

// Set sets the point´s underlying value.
func (t *t_Float32) Set(v float32) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Float32) Get() float32 { return t.data }

// ****************************************************************************

// Float64 represents the sunspec type float64.
type Float64 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v float64) error
	// Get returns the point´s underlying value.
	Get() float64
}

type t_Float64 struct {
	point
	data float64
}

var _ Float64 = (*t_Float64)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Float64) Valid() bool { return t.Get() != 0x7FC00000 }

// String formats the point´s value as string.
func (t *t_Float64) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Float64) Quantity() uint16 { return 4 }

// Encode puts the point´s value into a buffer.
func (t *t_Float64) Encode(buf []byte) error {
	binary.BigEndian.PutUint64(buf, math.Float64bits(t.Get()))
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Float64) Decode(buf []byte) error {
	return t.Set(math.Float64frombits(binary.BigEndian.Uint64(buf)))
}

// Set sets the point´s underlying value.
func (t *t_Float64) Set(v float64) error {
	t.data = v
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Float64) Get() float64 { return t.data }

// ****************************************************************************

// Ipaddr represents the sunspec type ipaddr.
type Ipaddr interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.IP) error
	// Get returns the point´s underlying value.
	Get() net.IP
	// Raw returns the point´s raw data.
	Raw() [4]byte
}

type t_Ipaddr struct {
	point
	data [4]byte
}

var _ Ipaddr = (*t_Ipaddr)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Ipaddr) Valid() bool { return t.data != [4]byte{} }

// String formats the point´s value as string.
func (t *t_Ipaddr) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Ipaddr) Quantity() uint16 { return uint16(len(t.data) / 2) }

// Encode puts the point´s value into a buffer.
func (t *t_Ipaddr) Encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Ipaddr) Decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *t_Ipaddr) Set(v net.IP) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Ipaddr) Get() net.IP { return append(net.IP(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *t_Ipaddr) Raw() (r [4]byte) {
	copy(r[:], t.Get())
	return r
}

// ****************************************************************************

// Ipaddr represents the sunspec type ipaddr.
type Ipv6addr interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.IP) error
	// Get returns the point´s underlying value.
	Get() net.IP
	// Raw returns the point´s raw data.
	Raw() [16]byte
}

type t_Ipv6addr struct {
	point
	data [16]byte
}

var _ Ipv6addr = (*t_Ipv6addr)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Ipv6addr) Valid() bool { return t.data != [16]byte{} }

// String formats the point´s value as string.
func (t *t_Ipv6addr) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Ipv6addr) Quantity() uint16 { return uint16(len(t.data) / 2) }

// Encode puts the point´s value into a buffer.
func (t *t_Ipv6addr) Encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Ipv6addr) Decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *t_Ipv6addr) Set(v net.IP) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Ipv6addr) Get() net.IP { return append(net.IP(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *t_Ipv6addr) Raw() (r [16]byte) {
	copy(r[:], t.Get())
	return r
}

// ****************************************************************************

// Eui48 represents the sunspec type eui48.
type Eui48 interface {
	// Point defines the generic behavior all sunspec types have in common.
	Point
	// Set sets the point´s underlying value.
	Set(v net.HardwareAddr) error
	// Get returns the point´s underlying value.
	Get() net.HardwareAddr
	// Raw returns the point´s raw data.
	Raw() [8]byte
}

type t_Eui48 struct {
	point
	data [8]byte
}

var _ Eui48 = (*t_Eui48)(nil)

// Valid specifies whether the underlying value is implemented by the device.
func (t *t_Eui48) Valid() bool { return true } //?

// String formats the point´s value as string.
func (t *t_Eui48) String() string { return fmt.Sprintf("%v", t.Get()) }

// Quantity returns the number of modbus registers required to store the underlying value.
func (t *t_Eui48) Quantity() uint16 { return uint16(len(t.data) / 2) }

// Encode puts the point´s value into a buffer.
func (t *t_Eui48) Encode(buf []byte) error {
	copy(buf, t.Get())
	return nil
}

// Decode sets the point´s value from a buffer.
func (t *t_Eui48) Decode(buf []byte) error {
	return t.Set(buf)
}

// Set sets the point´s underlying value.
func (t *t_Eui48) Set(v net.HardwareAddr) error {
	copy(t.data[:len(t.data)], v)
	return nil
}

// Get returns the point´s underlying value.
func (t *t_Eui48) Get() net.HardwareAddr { return append(net.HardwareAddr(nil), t.data[:]...) }

// Raw returns the point´s raw data.
func (t *t_Eui48) Raw() (r [8]byte) {
	copy(r[:], t.Get())
	return r
}
