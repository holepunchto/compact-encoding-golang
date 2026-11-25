package compactencoding

import (
	"errors"
	"fmt"
	"reflect"
)

// A Marshaler encodes itself, such as to pack fields into flags. Marshal
// defers to it for the value and for any field or element of the type.
type Marshaler interface {
	Preencode(state *State)
	Encode(state *State) error
}

// An Unmarshaler decodes itself. Unmarshal defers to it as Marshal defers to
// a Marshaler.
type Unmarshaler interface {
	Decode(state *State) error
}

// Marshal encodes v, which is encoded by its type:
//
//   - A Marshaler encodes itself.
//   - A struct encodes its exported fields in order, skipping any tagged
//     `compact:"-"`.
//   - A []byte encodes as Buffer and any other slice as Array.
//   - A bool, string, float, or integer encodes as the codec named after its
//     type, such as Uint16 for a uint16.
//   - A pointer encodes as the value it points to, and must not be nil.
func Marshal(v any) ([]byte, error) {
	rv, err := marshalValue(v)
	if err != nil {
		return nil, err
	}

	var state State

	if err := preencodeValue(&state, rv); err != nil {
		return nil, err
	}

	state.Allocate()

	if err := encodeValue(&state, rv); err != nil {
		return nil, err
	}

	return state.Buffer, nil
}

// Unmarshal decodes data into the value that v points to, as Marshal encodes
// it. Pointers are allocated as needed.
func Unmarshal(data []byte, v any) error {
	return DecodeFrom(&State{End: len(data), Buffer: data}, v)
}

// PreencodeInto adds the encoded size of v to state, as Marshal encodes it.
// Together with EncodeInto, it encodes several values into one buffer.
func PreencodeInto(state *State, v any) error {
	rv, err := marshalValue(v)
	if err != nil {
		return err
	}

	return preencodeValue(state, rv)
}

// EncodeInto encodes v at state.Start, as Marshal encodes it, once the buffer
// is allocated.
func EncodeInto(state *State, v any) error {
	rv, err := marshalValue(v)
	if err != nil {
		return err
	}

	return encodeValue(state, rv)
}

// DecodeFrom decodes a value at state.Start into the value that v points to,
// as Unmarshal decodes it.
func DecodeFrom(state *State, v any) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("compactencoding: decoding requires a non-nil pointer")
	}

	return decodeValue(state, rv.Elem())
}

var (
	errNilPointer = errors.New("compactencoding: nil pointer")

	marshalerType   = reflect.TypeFor[Marshaler]()
	unmarshalerType = reflect.TypeFor[Unmarshaler]()
	bytesType       = reflect.TypeFor[[]byte]()
)

// Copies a value that can't be addressed, so that the methods of its pointer,
// such as those of a Marshaler, can be called.
func marshalValue(v any) (reflect.Value, error) {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return rv, errNilPointer
	}

	if rv.Kind() == reflect.Pointer {
		return rv, nil
	}

	p := reflect.New(rv.Type())
	p.Elem().Set(rv)

	return p.Elem(), nil
}

func marshaler(rv reflect.Value) (Marshaler, bool) {
	if rv.CanAddr() && reflect.PointerTo(rv.Type()).Implements(marshalerType) {
		return rv.Addr().Interface().(Marshaler), true
	}

	if rv.Type().Implements(marshalerType) {
		return rv.Interface().(Marshaler), true
	}

	return nil, false
}

func unmarshaler(rv reflect.Value) (Unmarshaler, bool) {
	if rv.CanAddr() && reflect.PointerTo(rv.Type()).Implements(unmarshalerType) {
		return rv.Addr().Interface().(Unmarshaler), true
	}

	if rv.Type().Implements(unmarshalerType) {
		return rv.Interface().(Unmarshaler), true
	}

	return nil, false
}

func preencodeValue(state *State, rv reflect.Value) error {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return errNilPointer
		}

		return preencodeValue(state, rv.Elem())
	}

	if m, ok := marshaler(rv); ok {
		m.Preencode(state)

		return nil
	}

	switch {
	case rv.Kind() == reflect.Struct:
		return eachField(rv, func(f reflect.Value) error {
			return preencodeValue(state, f)
		})
	case isArray(rv.Type()):
		state.End += varintSize(uint64(rv.Len()))

		return eachElement(rv, func(e reflect.Value) error {
			return preencodeValue(state, e)
		})
	}

	c, err := scalarFor(rv.Type())
	if err != nil {
		return err
	}

	c.preencode(state, rv)

	return nil
}

func encodeValue(state *State, rv reflect.Value) error {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return errNilPointer
		}

		return encodeValue(state, rv.Elem())
	}

	if m, ok := marshaler(rv); ok {
		return m.Encode(state)
	}

	switch {
	case rv.Kind() == reflect.Struct:
		return eachField(rv, func(f reflect.Value) error {
			return encodeValue(state, f)
		})
	case isArray(rv.Type()):
		if err := encodeVarint(state, uint64(rv.Len())); err != nil {
			return err
		}

		return eachElement(rv, func(e reflect.Value) error {
			return encodeValue(state, e)
		})
	}

	c, err := scalarFor(rv.Type())
	if err != nil {
		return err
	}

	return c.encode(state, rv)
}

func decodeValue(state *State, rv reflect.Value) error {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		return decodeValue(state, rv.Elem())
	}

	if u, ok := unmarshaler(rv); ok {
		return u.Decode(state)
	}

	switch {
	case rv.Kind() == reflect.Struct:
		return eachField(rv, func(f reflect.Value) error {
			return decodeValue(state, f)
		})
	case isArray(rv.Type()):
		n, err := decodeVarint(state)
		if err != nil {
			return err
		}

		if n > MaxArrayLength {
			return ErrArrayTooBig
		}

		rv.Set(reflect.MakeSlice(rv.Type(), int(n), int(n)))

		return eachElement(rv, func(e reflect.Value) error {
			return decodeValue(state, e)
		})
	}

	c, err := scalarFor(rv.Type())
	if err != nil {
		return err
	}

	return c.decode(state, rv)
}

func isArray(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && !t.ConvertibleTo(bytesType)
}

func eachField(rv reflect.Value, fn func(reflect.Value) error) error {
	t := rv.Type()

	for i := range t.NumField() {
		f := t.Field(i)

		if !f.IsExported() || f.Tag.Get("compact") == "-" {
			continue
		}

		if err := fn(rv.Field(i)); err != nil {
			return fmt.Errorf("%s: %w", f.Name, err)
		}
	}

	return nil
}

func eachElement(rv reflect.Value, fn func(reflect.Value) error) error {
	for i := range rv.Len() {
		if err := fn(rv.Index(i)); err != nil {
			return fmt.Errorf("[%d]: %w", i, err)
		}
	}

	return nil
}

type scalarCodec struct {
	preencode func(*State, reflect.Value)
	encode    func(*State, reflect.Value) error
	decode    func(*State, reflect.Value) error
}

// Values are converted to and from the type of the codec, so that named
// types, such as `type ID uint32`, encode as their underlying type.
func scalar[T any](codec Codec[T]) scalarCodec {
	t := reflect.TypeFor[T]()

	return scalarCodec{
		preencode: func(state *State, rv reflect.Value) {
			codec.Preencode(state, rv.Convert(t).Interface().(T))
		},
		encode: func(state *State, rv reflect.Value) error {
			return codec.Encode(state, rv.Convert(t).Interface().(T))
		},
		decode: func(state *State, rv reflect.Value) error {
			v, err := codec.Decode(state)
			if err != nil {
				return err
			}

			rv.Set(reflect.ValueOf(v).Convert(rv.Type()))

			return nil
		},
	}
}

var scalars = map[reflect.Kind]scalarCodec{
	reflect.Bool:    scalar(Bool),
	reflect.String:  scalar(String),
	reflect.Uint:    scalar(Uint),
	reflect.Uint8:   scalar(Uint8),
	reflect.Uint16:  scalar(Uint16),
	reflect.Uint32:  scalar(Uint32),
	reflect.Uint64:  scalar(Uint64),
	reflect.Int:     scalar(Int),
	reflect.Int8:    scalar(Int8),
	reflect.Int16:   scalar(Int16),
	reflect.Int32:   scalar(Int32),
	reflect.Int64:   scalar(Int64),
	reflect.Float32: scalar(Float32),
	reflect.Float64: scalar(Float64),
	reflect.Slice:   scalar(Buffer),
}

func scalarFor(t reflect.Type) (scalarCodec, error) {
	c, ok := scalars[t.Kind()]

	if !ok {
		return c, fmt.Errorf("compactencoding: unsupported type %s", t)
	}

	return c, nil
}
