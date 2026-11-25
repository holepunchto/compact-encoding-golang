package compactencoding_test

import (
	"encoding/hex"
	"errors"
	"math"
	"reflect"
	"testing"

	c "github.com/holepunchto/compact-encoding-golang/v2"
)

func encodes[T any](t *testing.T, codec c.Codec[T], value T, want string) {
	t.Helper()

	b, err := c.Encode(codec, value)
	if err != nil {
		t.Fatalf("encode %v: %v", value, err)
	}

	if got := hex.EncodeToString(b); got != want {
		t.Errorf("encode %v = %s, want %s", value, got, want)
	}

	decoded, err := c.Decode(codec, b)
	if err != nil {
		t.Fatalf("decode %s: %v", want, err)
	}

	if !reflect.DeepEqual(decoded, value) {
		t.Errorf("decode %s = %v, want %v", want, decoded, value)
	}
}

func decodes[T any](t *testing.T, codec c.Codec[T], input string, want T, read int) {
	t.Helper()

	b, _ := hex.DecodeString(input)
	state := &c.State{End: len(b), Buffer: b}

	got, err := codec.Decode(state)
	if err != nil {
		t.Fatalf("decode %s: %v", input, err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("decode %s = %#v, want %#v", input, got, want)
	}

	if state.Start != read {
		t.Errorf("decode %s read %d bytes, want %d", input, state.Start, read)
	}
}

func rejects[T any](t *testing.T, codec c.Codec[T], input string, want error) {
	t.Helper()

	b, _ := hex.DecodeString(input)

	if _, err := c.Decode(codec, b); !errors.Is(err, want) {
		t.Errorf("decode %s = %v, want %v", input, err, want)
	}
}

func TestUint(t *testing.T) {
	encodes(t, c.Uint, 0, "00")
	encodes(t, c.Uint, 252, "fc")
	encodes(t, c.Uint, 253, "fdfd00")
	encodes(t, c.Uint, 65535, "fdffff")
	encodes(t, c.Uint, 65536, "fe00000100")
	encodes(t, c.Uint, 0xffffffff, "feffffffff")
	if v := uint64(1) << 32; v <= math.MaxUint {
		encodes(t, c.Uint, uint(v), "ff0000000001000000")
	} else {
		rejects(t, c.Uint, "ff0000000001000000", c.ErrOverflow)
	}

	decodes(t, c.Uint, "fd0100", 1, 3)
	decodes(t, c.Uint, "0161", 1, 1)

	rejects(t, c.Uint, "fd01", c.ErrOutOfBounds)
	rejects(t, c.Uint, "", c.ErrOutOfBounds)
}

func TestInt(t *testing.T) {
	encodes(t, c.Int, 0, "00")
	encodes(t, c.Int, -1, "01")
	encodes(t, c.Int, 1, "02")
	encodes(t, c.Int, -2, "03")
	encodes(t, c.Int, -64, "7f")
	encodes(t, c.Int, 127, "fdfe00")
	encodes(t, c.Int, -128, "fdff00")
	encodes(t, c.Int, math.MaxInt32, "fefeffffff")
	encodes(t, c.Int, math.MinInt32, "feffffffff")
}

func TestFixedWidth(t *testing.T) {
	encodes(t, c.Uint8, 255, "ff")
	encodes(t, c.Uint16, 258, "0201")
	encodes(t, c.Uint32, 0x01020304, "04030201")
	encodes(t, c.Uint64, 0x0001020304050607, "0706050403020100")
	encodes(t, c.Int8, -128, "ff")
	encodes(t, c.Int8, 127, "fe")
	encodes(t, c.Int16, -2, "0300")
	encodes(t, c.Int32, -1, "01000000")
	encodes(t, c.Int64, -3, "0500000000000000")

	rejects(t, c.Uint32, "010203", c.ErrOutOfBounds)
}

func TestFloat(t *testing.T) {
	encodes(t, c.Float32, 0.5, "0000003f")
	encodes(t, c.Float64, 1, "000000000000f03f")
	encodes(t, c.Float64, math.Inf(-1), "000000000000f0ff")
}

func TestBool(t *testing.T) {
	encodes(t, c.Bool, true, "01")
	encodes(t, c.Bool, false, "00")

	decodes(t, c.Bool, "02", false, 1)

	rejects(t, c.Bool, "", c.ErrOutOfBounds)
}

func TestBuffer(t *testing.T) {
	encodes(t, c.Buffer, []byte("hi"), "026869")
	encodes(t, c.Buffer, []byte{}, "00")

	decodes(t, c.Buffer, "00", []byte{}, 1)
	decodes(t, c.OptionalBuffer, "00", nil, 1)
	decodes(t, c.OptionalBuffer, "0161", []byte("a"), 2)

	rejects(t, c.Buffer, "0201", c.ErrOutOfBounds)
	rejects(t, c.Buffer, "ffffffffffffffffff", c.ErrOutOfBounds)

	if b, _ := c.Encode(c.Buffer, nil); hex.EncodeToString(b) != "00" {
		t.Errorf("encode nil = %x, want 00", b)
	}
}

func TestString(t *testing.T) {
	encodes(t, c.String, "abc", "03616263")
	encodes(t, c.String, "æøå", "06c3a6c3b8c3a5")
	encodes(t, c.String, "😀", "04f09f9880")

	decodes(t, c.String, "03e28241", "�A", 4)
	decodes(t, c.String, "01ff", "�", 2)
	decodes(t, c.String, "03eda080", "���", 4)
	decodes(t, c.String, "04f0908041", "�A", 5)
	decodes(t, c.String, "02c0af", "��", 3)
	decodes(t, c.String, "03e0809f", "���", 4)

	if b, _ := c.Encode(c.String, "\xe2\x82A"); hex.EncodeToString(b) != "04efbfbd41" {
		t.Errorf("encode ill-formed string = %x, want 04efbfbd41", b)
	}
}

func TestArray(t *testing.T) {
	encodes(t, c.Array(c.Uint), []uint{1, 253}, "0201fdfd00")
	encodes(t, c.Array(c.String), []string{}, "00")

	rejects(t, c.Array(c.Uint), "fe01001000", c.ErrArrayTooBig)
	rejects(t, c.Array(c.Uint), "0201", c.ErrOutOfBounds)
}

func TestFrame(t *testing.T) {
	encodes(t, c.Frame(c.Uint), 253, "03fdfd00")

	decodes(t, c.Frame(c.Uint), "0201ff07", 1, 3)

	rejects(t, c.Frame(c.Uint), "0301", c.ErrOutOfBounds)
	rejects(t, c.Frame(c.Uint32), "0201ff", c.ErrOutOfBounds)
}

func TestRecord(t *testing.T) {
	encodes(t, c.Record(c.String, c.Uint), map[string]uint{"b": 2, "a": 1}, "02016101016202")

	decodes(t, c.Record(c.String, c.Uint), "02016101016102", map[string]uint{"a": 2}, 7)

	rejects(t, c.Record(c.String, c.Uint), "ff0000000000000001", c.ErrOutOfBounds)
}

func TestJSON(t *testing.T) {
	type value struct {
		A string `json:"a"`
		B []int  `json:"b"`
	}

	encodes(t, c.JSON[value](), value{A: "<&>", B: []int{1, 2}}, "157b2261223a223c263e222c2262223a5b312c325d7d")

	if _, err := c.Encode(c.JSON[func()](), nil); err == nil {
		t.Error("encode of a value JSON can't hold succeeded")
	}
}
