package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInt(t *testing.T) {
	state := &State{}
	Encoder := NewInt()

	Convey("int can be Encoded", t, func() {
		Encoder.Preencode(state, 42)
		So(state, ShouldResemble, &State{Start: 0, End: 1, Buffer: nil})
		Encoder.Preencode(state, -4200)
		So(state, ShouldResemble, &State{Start: 0, End: 4, Buffer: nil})

		state.Allocate()

		err := Encoder.Encode(state, 42)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 1, End: 4, Buffer: []byte{84, 0, 0, 0}})

		err = Encoder.Encode(state, -4200)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 4, End: 4, Buffer: []byte{84, 0xfd, 207, 32}})

		state.Rewind()

		value, err := Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 42)

		value, err = Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, -4200)

		_, err = Encoder.Decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
