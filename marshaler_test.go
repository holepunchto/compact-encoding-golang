package compactencoding

import (
	"bytes"
	"testing"
)

// flagged packs its bool into a flags byte, the way hyperschema does, which reflection alone cannot express.
type flagged struct {
	ID    uint
	Happy bool
	Name  string
}

func (f *flagged) Preencode(state *State) {
	NewUint().Preencode(state, f.ID)
	state.End++
	if f.Name != "" {
		NewString().Preencode(state, f.Name)
	}
}

func (f *flagged) Encode(state *State) error {
	if err := NewUint().Encode(state, f.ID); err != nil {
		return err
	}
	var flags uint
	if f.Happy {
		flags |= 1
	}
	if f.Name != "" {
		flags |= 2
	}
	if err := NewUint().Encode(state, flags); err != nil {
		return err
	}
	if f.Name != "" {
		return NewString().Encode(state, f.Name)
	}
	return nil
}

func (f *flagged) Decode(state *State) error {
	var err error
	if f.ID, err = NewUint().Decode(state); err != nil {
		return err
	}
	flags, err := NewUint().Decode(state)
	if err != nil {
		return err
	}
	f.Happy = flags&1 != 0
	if flags&2 != 0 {
		if f.Name, err = NewString().Decode(state); err != nil {
			return err
		}
	}
	return nil
}

type holder struct {
	Items []flagged
	One   flagged
	Ptr   *flagged
}

func TestMarshalerByValueAndPointer(t *testing.T) {
	want := []byte{7, 3, 2, 'h', 'i'}
	for _, v := range []any{flagged{ID: 7, Happy: true, Name: "hi"}, &flagged{ID: 7, Happy: true, Name: "hi"}} {
		got, err := Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("got %x, want %x", got, want)
		}
	}

	var decoded flagged
	if err := Unmarshal(want, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != 7 || !decoded.Happy || decoded.Name != "hi" {
		t.Fatalf("decoded %+v", decoded)
	}
}

func TestMarshalerNested(t *testing.T) {
	in := holder{
		Items: []flagged{{ID: 1}, {ID: 2, Happy: true}},
		One:   flagged{ID: 3, Name: "x"},
		Ptr:   &flagged{ID: 4},
	}
	data, err := Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{2, 1, 0, 2, 1, 3, 2, 1, 'x', 4, 0}
	if !bytes.Equal(data, want) {
		t.Fatalf("got %x, want %x", data, want)
	}

	var out holder
	if err := Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if len(out.Items) != 2 || !out.Items[1].Happy || out.One.Name != "x" || out.Ptr == nil || out.Ptr.ID != 4 {
		t.Fatalf("decoded %+v", out)
	}
}
