package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalScalars(t *testing.T) {
	Convey("Marshal/Unmarshal scalar types", t, func() {
		Convey("string", func() {
			data, err := Marshal("hello")
			So(err, ShouldBeNil)
			var got string
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, "hello")
		})

		Convey("bool true", func() {
			data, err := Marshal(true)
			So(err, ShouldBeNil)
			var got bool
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldBeTrue)
		})

		Convey("bool false", func() {
			data, err := Marshal(false)
			So(err, ShouldBeNil)
			var got bool
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldBeFalse)
		})

		Convey("uint8 (fixed)", func() {
			data, err := Marshal(uint8(255))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{255})
			var got uint8
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 255)
		})

		Convey("uint16 (fixed LE)", func() {
			data, err := Marshal(uint16(0x0102))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{0x02, 0x01})
			var got uint16
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 0x0102)
		})

		Convey("uint32 (fixed LE)", func() {
			data, err := Marshal(uint32(0x01020304))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{0x04, 0x03, 0x02, 0x01})
			var got uint32
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 0x01020304)
		})

		Convey("uint64 (fixed LE)", func() {
			data, err := Marshal(uint64(0x0102030405060708))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01})
			var got uint64
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 0x0102030405060708)
		})

		Convey("uint (varuint, small)", func() {
			data, err := Marshal(uint(42))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{42})
			var got uint
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 42)
		})

		Convey("uint (varuint, large)", func() {
			data, err := Marshal(uint(0x1234))
			So(err, ShouldBeNil)
			So(data, ShouldResemble, []byte{0xfd, 0x34, 0x12})
			var got uint
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 0x1234)
		})

		Convey("int (varint, zigzag positive)", func() {
			data, err := Marshal(42)
			So(err, ShouldBeNil)
			var got int
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, 42)
		})

		Convey("int (varint, zigzag negative)", func() {
			data, err := Marshal(-1)
			So(err, ShouldBeNil)
			// -1 zigzag → 1 → varuint 1 byte
			So(data, ShouldResemble, []byte{1})
			var got int
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldEqual, -1)
		})

		Convey("[]byte (buffer)", func() {
			data, err := Marshal([]byte{1, 2, 3})
			So(err, ShouldBeNil)
			var got []byte
			So(Unmarshal(data, &got), ShouldBeNil)
			So(got, ShouldResemble, []byte{1, 2, 3})
		})

		Convey("nil []byte encodes as zero-length buffer", func() {
			data, err := Marshal([]byte(nil))
			So(err, ShouldBeNil)
			var got []byte
			So(Unmarshal(data, &got), ShouldBeNil)
			So(len(got), ShouldEqual, 0)
		})
	})
}

func TestMarshalStruct(t *testing.T) {
	type ChatRequest struct {
		SessionID string
		Message   string
		Model     string
	}

	Convey("Marshal/Unmarshal simple struct", t, func() {
		original := ChatRequest{
			SessionID: "sess-abc",
			Message:   "hello world",
			Model:     "claude-3-5-sonnet",
		}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		var got ChatRequest
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got, ShouldResemble, original)
	})

	Convey("bytes match manual encoding", t, func() {
		// Verify Marshal produces the same bytes as hand-coded Encode calls.
		original := ChatRequest{SessionID: "s", Message: "m", Model: "x"}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		state := NewState()
		enc := NewString()
		enc.Preencode(state, "s")
		enc.Preencode(state, "m")
		enc.Preencode(state, "x")
		state.Allocate()
		_ = enc.Encode(state, "s")
		_ = enc.Encode(state, "m")
		_ = enc.Encode(state, "x")

		So(data, ShouldResemble, state.Buffer)
	})

	Convey("pointer receiver works identically", t, func() {
		original := &ChatRequest{SessionID: "p", Message: "q", Model: "r"}
		data, err := Marshal(original)
		So(err, ShouldBeNil)
		var got ChatRequest
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got, ShouldResemble, *original)
	})
}

func TestMarshalMixedTypes(t *testing.T) {
	type Msg struct {
		Type    uint8
		Seq     uint
		Payload []byte
		Label   string
		Ok      bool
	}

	Convey("round-trip struct with mixed scalar types", t, func() {
		original := Msg{
			Type:    3,
			Seq:     1000,
			Payload: []byte{0xde, 0xad, 0xbe, 0xef},
			Label:   "test",
			Ok:      true,
		}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		var got Msg
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got, ShouldResemble, original)
	})
}

func TestMarshalNestedStruct(t *testing.T) {
	type Scope struct {
		AgentID string
		Channel string
	}
	type SessionRequest struct {
		Scope   Scope
		Peer    string
		Version uint8
	}

	Convey("round-trip nested struct", t, func() {
		original := SessionRequest{
			Scope:   Scope{AgentID: "agent-1", Channel: "cli"},
			Peer:    "user-42",
			Version: 1,
		}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		var got SessionRequest
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got, ShouldResemble, original)
	})
}

func TestMarshalSliceField(t *testing.T) {
	type Batch struct {
		Tags  []string
		Count uint32
	}

	Convey("round-trip []string slice field", t, func() {
		original := Batch{
			Tags:  []string{"alpha", "beta", "gamma"},
			Count: 3,
		}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		var got Batch
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got, ShouldResemble, original)
	})

	Convey("nil slice encodes as length-zero array", t, func() {
		original := Batch{Tags: nil, Count: 0}
		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		var got Batch
		So(Unmarshal(data, &got), ShouldBeNil)
		So(len(got.Tags), ShouldEqual, 0)
		So(got.Count, ShouldEqual, 0)
	})
}

func TestMarshalSkipTag(t *testing.T) {
	type WithSkip struct {
		Name    string
		Secret  string `compact:"-"`
		Version uint8
	}

	Convey("compact:\"-\" fields are omitted from the wire", t, func() {
		original := WithSkip{Name: "foo", Secret: "do-not-encode", Version: 2}

		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		// Should match encoding of just Name + Version.
		state := NewState()
		NewString().Preencode(state, "foo")
		NewUint8().Preencode(state, uint8(2))
		state.Allocate()
		_ = NewString().Encode(state, "foo")
		_ = NewUint8().Encode(state, 2)
		So(data, ShouldResemble, state.Buffer)

		// Decode restores Name and Version; Secret stays zero-value.
		var got WithSkip
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got.Name, ShouldEqual, "foo")
		So(got.Secret, ShouldEqual, "")
		So(got.Version, ShouldEqual, 2)
	})
}

func TestMarshalUnexportedFieldsIgnored(t *testing.T) {
	type WithPrivate struct {
		Public  string
		private string //nolint:unused
	}

	Convey("unexported fields are silently skipped", t, func() {
		original := WithPrivate{Public: "visible"}
		data, err := Marshal(&original)
		So(err, ShouldBeNil)

		// Bytes should equal encoding of just the Public string.
		expected, _ := Marshal("visible")
		So(data, ShouldResemble, expected)

		var got WithPrivate
		So(Unmarshal(data, &got), ShouldBeNil)
		So(got.Public, ShouldEqual, "visible")
	})
}

func TestMarshalErrors(t *testing.T) {
	Convey("Marshal errors", t, func() {
		Convey("nil pointer returns error", func() {
			var p *struct{ X string }
			_, err := Marshal(p)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "nil pointer")
		})

		Convey("Unmarshal into non-pointer returns error", func() {
			var s struct{ X string }
			err := Unmarshal([]byte{}, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "non-nil pointer")
		})
	})
}
