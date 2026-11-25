package compactencoding_test

import (
	"encoding/hex"
	"errors"
	"reflect"
	"testing"

	c "github.com/holepunchto/compact-encoding-golang/v2"
)

type ID uint32

type item struct {
	ID      ID
	Name    string
	Score   float64
	Tags    []string
	Data    []byte
	Skipped string `compact:"-"`
	hidden  string
}

func TestMarshalStruct(t *testing.T) {
	in := item{ID: 7, Name: "a", Score: 1, Tags: []string{"x"}, Data: []byte{0xff}, Skipped: "s", hidden: "h"}

	b, err := c.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	want := "07000000" + "0161" + "000000000000f03f" + "010178" + "01ff"

	if got := hex.EncodeToString(b); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	var out item

	if err := c.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}

	in.Skipped, in.hidden = "", ""

	if !reflect.DeepEqual(out, in) {
		t.Fatalf("got %+v, want %+v", out, in)
	}
}

// Packs its bool into a flags byte, which reflection alone can't express.
type flagged struct {
	ID    uint
	Happy bool
	Name  string
}

func (f *flagged) Preencode(state *c.State) {
	c.Uint.Preencode(state, f.ID)
	c.Uint.Preencode(state, 0)

	if f.Name != "" {
		c.String.Preencode(state, f.Name)
	}
}

func (f *flagged) Encode(state *c.State) error {
	var flags uint

	if f.Happy {
		flags |= 1
	}

	if f.Name != "" {
		flags |= 2
	}

	if err := c.Uint.Encode(state, f.ID); err != nil {
		return err
	}

	if err := c.Uint.Encode(state, flags); err != nil {
		return err
	}

	if f.Name != "" {
		return c.String.Encode(state, f.Name)
	}

	return nil
}

func (f *flagged) Decode(state *c.State) error {
	var err error

	if f.ID, err = c.Uint.Decode(state); err != nil {
		return err
	}

	flags, err := c.Uint.Decode(state)
	if err != nil {
		return err
	}

	f.Happy = flags&1 != 0

	if flags&2 != 0 {
		f.Name, err = c.String.Decode(state)
	}

	return err
}

func TestMarshalDefersToMarshaler(t *testing.T) {
	type holder struct {
		Items []flagged
		One   flagged
		Ptr   *flagged
	}

	in := holder{
		Items: []flagged{{ID: 1}, {ID: 2, Happy: true}},
		One:   flagged{ID: 3, Name: "x"},
		Ptr:   &flagged{ID: 4},
	}

	for _, v := range []any{in, &in} {
		b, err := c.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := hex.EncodeToString(b), "02010002010302017804"+"00"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		var out holder

		if err := c.Unmarshal(b, &out); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(out, in) {
			t.Fatalf("got %+v, want %+v", out, in)
		}
	}
}

func TestMarshalErrors(t *testing.T) {
	type pointer struct{ P *uint }

	if _, err := c.Marshal(pointer{}); err == nil {
		t.Error("marshal of a nil pointer succeeded")
	}

	if _, err := c.Marshal(map[string]int{}); err == nil {
		t.Error("marshal of a map succeeded")
	}

	var n uint

	if err := c.Unmarshal([]byte{1}, n); err == nil {
		t.Error("unmarshal into a non-pointer succeeded")
	}

	var s []uint

	if err := c.Unmarshal([]byte{0xfe, 0x01, 0x00, 0x10, 0x00}, &s); !errors.Is(err, c.ErrArrayTooBig) {
		t.Errorf("unmarshal of a huge array = %v, want ErrArrayTooBig", err)
	}
}

func TestEncodeIntoSharesABuffer(t *testing.T) {
	var state c.State

	for _, v := range []any{uint(1), "a"} {
		if err := c.PreencodeInto(&state, v); err != nil {
			t.Fatal(err)
		}
	}

	state.Allocate()

	for _, v := range []any{uint(1), "a"} {
		if err := c.EncodeInto(&state, v); err != nil {
			t.Fatal(err)
		}
	}

	state.Start = 0

	var n uint
	var s string

	if err := c.DecodeFrom(&state, &n); err != nil {
		t.Fatal(err)
	}

	if err := c.DecodeFrom(&state, &s); err != nil {
		t.Fatal(err)
	}

	if n != 1 || s != "a" {
		t.Fatalf("got %d, %q", n, s)
	}
}
