package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArray(t *testing.T) {
	Convey("uint8 array", t, func() {
		state := &State{}
		Encoder := NewArray(NewUint())
		Encoder.Preencode(state, []uint{1, 2, 3})
		state.Allocate()

		err := Encoder.Encode(state, []uint{1, 2, 3})

		So(err, ShouldBeNil)
		So(state.Buffer, ShouldEqual, []byte{3, 1, 2, 3})
		state.Rewind()
		value, err := Encoder.Decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []uint{1, 2, 3})
	})

	Convey("string array", t, func() {
		state := &State{}
		Encoder := NewArray(NewString())
		Encoder.Preencode(state, []string{"hello", "world"})
		state.Allocate()

		err := Encoder.Encode(state, []string{"hello", "world"})

		So(err, ShouldBeNil)
		So(state.Buffer, ShouldEqual, []byte{2, 5, 104, 101, 108, 108, 111, 5, 119, 111, 114, 108, 100})
		state.Rewind()
		value, err := Encoder.Decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []string{"hello", "world"})
	})
}

func TestArrayExisting(t *testing.T) {
	Convey("existing buffer", t, func() {
		state := &State{}
		state.Buffer = []byte{2, 3, 111, 110, 101, 3, 116, 119, 111}
		state.End = 9

		Encoder := NewArray(NewString())
		value, err := Encoder.Decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []string{"one", "two"})
	})
}
