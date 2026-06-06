package compactencoding

import (
	"fmt"
	"reflect"
)

// Marshal encodes v into compact-encoding bytes. v may be a struct, pointer to
// struct, or any scalar type supported by the library. Struct fields are
// encoded in declaration order. Unexported fields and fields tagged
// `compact:"-"` are skipped.
//
// Go type → codec mapping:
//
//	string        → string (varuint length + UTF-8 bytes)
//	bool          → bool (1 byte)
//	uint8         → uint8 (1 byte, fixed)
//	uint16        → uint16 (2 bytes LE, fixed)
//	uint32        → uint32 (4 bytes LE, fixed)
//	uint64        → uint64 (8 bytes LE, fixed)
//	uint          → uint (varuint, 1–9 bytes)
//	int8          → int8 (1 byte, zigzag)
//	int16         → int16 (2 bytes LE, zigzag)
//	int32         → int32 (4 bytes LE, zigzag)
//	int64         → int64 (8 bytes LE, zigzag)
//	int           → int (varint, zigzag)
//	[]byte        → buffer (varuint length + raw bytes)
//	[]T           → array (varuint length + encoded elements)
//	struct        → fields encoded in declaration order
func Marshal(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("compact: Marshal of nil pointer")
		}
		rv = rv.Elem()
	}
	state := NewState()
	if err := preencodeReflect(state, rv); err != nil {
		return nil, err
	}
	state.Allocate()
	if err := encodeReflect(state, rv); err != nil {
		return nil, err
	}
	return state.Buffer, nil
}

// Unmarshal decodes compact-encoding bytes into the value pointed to by v.
// v must be a non-nil pointer. Fields are decoded in the same order Marshal
// encodes them; the struct layout of v must match the encoded layout exactly.
func Unmarshal(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("compact: Unmarshal requires a non-nil pointer")
	}
	state := &State{
		End:    uint(len(data)),
		Buffer: data,
	}
	return decodeReflect(state, rv.Elem())
}

// PreencodeInto accumulates the encoded size of v into state.End without
// writing any bytes. Use together with EncodeInto to pack multiple values into
// a single pre-allocated buffer.
func PreencodeInto(state *State, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("compact: PreencodeInto of nil pointer")
		}
		rv = rv.Elem()
	}
	return preencodeReflect(state, rv)
}

// EncodeInto writes v into state.Buffer at state.Start, advancing state.Start.
// state.Buffer must be allocated (e.g. via state.Allocate) before calling.
func EncodeInto(state *State, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("compact: EncodeInto of nil pointer")
		}
		rv = rv.Elem()
	}
	return encodeReflect(state, rv)
}

// DecodeFrom reads from state into v, advancing state.Start past the decoded bytes.
// v must be a non-nil pointer. Use when decoding multiple values from a shared
// State created over an existing byte slice.
func DecodeFrom(state *State, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("compact: DecodeFrom requires a non-nil pointer")
	}
	return decodeReflect(state, rv.Elem())
}


func preencodeReflect(state *State, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() || f.Tag.Get("compact") == "-" {
				continue
			}
			if err := preencodeReflect(state, rv.Field(i)); err != nil {
				return fmt.Errorf("%s: %w", f.Name, err)
			}
		}
		return nil

	case reflect.Ptr:
		if rv.IsNil() {
			return fmt.Errorf("compact: nil pointer in value")
		}
		return preencodeReflect(state, rv.Elem())

	case reflect.String:
		NewString().Preencode(state, rv.String())
	case reflect.Bool:
		NewBool().Preencode(state, rv.Bool())
	case reflect.Uint8:
		NewUint8().Preencode(state, uint8(rv.Uint()))
	case reflect.Uint16:
		NewUint16().Preencode(state, uint16(rv.Uint()))
	case reflect.Uint32:
		NewUint32().Preencode(state, uint32(rv.Uint()))
	case reflect.Uint64:
		NewUint64().Preencode(state, uint64(rv.Uint()))
	case reflect.Uint:
		NewUint().Preencode(state, uint(rv.Uint()))
	case reflect.Int8:
		NewInt8().Preencode(state, int8(rv.Int()))
	case reflect.Int16:
		NewInt16().Preencode(state, int16(rv.Int()))
	case reflect.Int32:
		NewInt32().Preencode(state, int32(rv.Int()))
	case reflect.Int64:
		NewInt64().Preencode(state, int64(rv.Int()))
	case reflect.Int:
		NewInt().Preencode(state, int(rv.Int()))

	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			if rv.IsNil() {
				NewBuffer().Preencode(state, nil)
			} else {
				NewBuffer().Preencode(state, rv.Bytes())
			}
			return nil
		}
		length := uint(0)
		if !rv.IsNil() {
			length = uint(rv.Len())
		}
		NewUint().Preencode(state, length)
		for i := 0; i < int(length); i++ {
			if err := preencodeReflect(state, rv.Index(i)); err != nil {
				return fmt.Errorf("[%d]: %w", i, err)
			}
		}
		return nil

	default:
		return fmt.Errorf("compact: unsupported type %s", rv.Type())
	}
	return nil
}

func encodeReflect(state *State, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() || f.Tag.Get("compact") == "-" {
				continue
			}
			if err := encodeReflect(state, rv.Field(i)); err != nil {
				return fmt.Errorf("%s: %w", f.Name, err)
			}
		}
		return nil

	case reflect.Ptr:
		if rv.IsNil() {
			return fmt.Errorf("compact: nil pointer in value")
		}
		return encodeReflect(state, rv.Elem())

	case reflect.String:
		return NewString().Encode(state, rv.String())
	case reflect.Bool:
		return NewBool().Encode(state, rv.Bool())
	case reflect.Uint8:
		return NewUint8().Encode(state, uint8(rv.Uint()))
	case reflect.Uint16:
		return NewUint16().Encode(state, uint16(rv.Uint()))
	case reflect.Uint32:
		return NewUint32().Encode(state, uint32(rv.Uint()))
	case reflect.Uint64:
		return NewUint64().Encode(state, uint64(rv.Uint()))
	case reflect.Uint:
		return NewUint().Encode(state, uint(rv.Uint()))
	case reflect.Int8:
		return NewInt8().Encode(state, int8(rv.Int()))
	case reflect.Int16:
		return NewInt16().Encode(state, int16(rv.Int()))
	case reflect.Int32:
		return NewInt32().Encode(state, int32(rv.Int()))
	case reflect.Int64:
		return NewInt64().Encode(state, int64(rv.Int()))
	case reflect.Int:
		return NewInt().Encode(state, int(rv.Int()))

	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			if rv.IsNil() {
				return NewBuffer().Encode(state, nil)
			}
			return NewBuffer().Encode(state, rv.Bytes())
		}
		length := uint(0)
		if !rv.IsNil() {
			length = uint(rv.Len())
		}
		if err := NewUint().Encode(state, length); err != nil {
			return err
		}
		for i := 0; i < int(length); i++ {
			if err := encodeReflect(state, rv.Index(i)); err != nil {
				return fmt.Errorf("[%d]: %w", i, err)
			}
		}
		return nil

	default:
		return fmt.Errorf("compact: unsupported type %s", rv.Type())
	}
}

func decodeReflect(state *State, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		t := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() || f.Tag.Get("compact") == "-" {
				continue
			}
			if err := decodeReflect(state, rv.Field(i)); err != nil {
				return fmt.Errorf("%s: %w", f.Name, err)
			}
		}
		return nil

	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return decodeReflect(state, rv.Elem())

	case reflect.String:
		v, err := NewString().Decode(state)
		if err != nil {
			return err
		}
		rv.SetString(v)

	case reflect.Bool:
		v, err := NewBool().Decode(state)
		if err != nil {
			return err
		}
		rv.SetBool(v)

	case reflect.Uint8:
		v, err := NewUint8().Decode(state)
		if err != nil {
			return err
		}
		rv.SetUint(uint64(v))

	case reflect.Uint16:
		v, err := NewUint16().Decode(state)
		if err != nil {
			return err
		}
		rv.SetUint(uint64(v))

	case reflect.Uint32:
		v, err := NewUint32().Decode(state)
		if err != nil {
			return err
		}
		rv.SetUint(uint64(v))

	case reflect.Uint64:
		v, err := NewUint64().Decode(state)
		if err != nil {
			return err
		}
		rv.SetUint(v)

	case reflect.Uint:
		v, err := NewUint().Decode(state)
		if err != nil {
			return err
		}
		rv.SetUint(uint64(v))

	case reflect.Int8:
		v, err := NewInt8().Decode(state)
		if err != nil {
			return err
		}
		rv.SetInt(int64(v))

	case reflect.Int16:
		v, err := NewInt16().Decode(state)
		if err != nil {
			return err
		}
		rv.SetInt(int64(v))

	case reflect.Int32:
		v, err := NewInt32().Decode(state)
		if err != nil {
			return err
		}
		rv.SetInt(int64(v))

	case reflect.Int64:
		v, err := NewInt64().Decode(state)
		if err != nil {
			return err
		}
		rv.SetInt(v)

	case reflect.Int:
		v, err := NewInt().Decode(state)
		if err != nil {
			return err
		}
		rv.SetInt(int64(v))

	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			v, err := NewBuffer().Decode(state)
			if err != nil {
				return err
			}
			rv.SetBytes(v)
			return nil
		}
		length, err := NewUint().Decode(state)
		if err != nil {
			return err
		}
		slice := reflect.MakeSlice(rv.Type(), int(length), int(length))
		for i := 0; i < int(length); i++ {
			if err := decodeReflect(state, slice.Index(i)); err != nil {
				return fmt.Errorf("[%d]: %w", i, err)
			}
		}
		rv.Set(slice)
		return nil

	default:
		return fmt.Errorf("compact: unsupported type %s", rv.Type())
	}
	return nil
}
