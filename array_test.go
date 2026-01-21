package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArray(t *testing.T) {
	Convey("uint8 array", t, func() {
		state := &State{}
		encoder := NewArray(NewUint())
		encoder.preencode(state, []uint{1, 2, 3})
		state.Allocate()

		err := encoder.encode(state, []uint{1, 2, 3})

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{3, 1, 2, 3})
		state.Rewind()
		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []uint{1, 2, 3})
	})

	Convey("string array", t, func() {
		state := &State{}
		encoder := NewArray(NewString())
		encoder.preencode(state, []string{"hello", "world"})
		state.Allocate()

		err := encoder.encode(state, []string{"hello", "world"})

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{2, 5, 104, 101, 108, 108, 111, 5, 119, 111, 114, 108, 100})
		state.Rewind()
		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []string{"hello", "world"})
	})
}
