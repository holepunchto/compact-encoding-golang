package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInt(t *testing.T) {
	state := &State{}
	encoder := NewInt()

	Convey("int can be encoded", t, func() {
		encoder.preencode(state, 42)
		So(state, ShouldResemble, &State{start: 0, end: 1, buffer: nil})
		encoder.preencode(state, -4200)
		So(state, ShouldResemble, &State{start: 0, end: 4, buffer: nil})

		state.Allocate()

		err := encoder.encode(state, 42)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 1, end: 4, buffer: []byte{84, 0, 0, 0}})

		err = encoder.encode(state, -4200)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 4, end: 4, buffer: []byte{84, 0xfd, 207, 32}})

		state.Rewind()

		value, err := encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 42)

		value, err = encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, -4200)

		_, err = encoder.decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
