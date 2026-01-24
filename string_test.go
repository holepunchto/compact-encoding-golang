package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	state := &State{}
	Encoder := NewString()

	Convey("string can be Encoded", t, func() {
		Encoder.Preencode(state, "🌾")
		So(state, ShouldResemble, &State{Start: 0, End: 5, Buffer: nil})
		Encoder.Preencode(state, "høsten er fin")
		So(state, ShouldResemble, &State{Start: 0, End: 20, Buffer: nil})

		state.Allocate()
		err := Encoder.Encode(state, "🌾")
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 5, End: 20, Buffer: []byte{4, 240, 159, 140, 190, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
		err = Encoder.Encode(state, "høsten er fin")
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 20, End: 20, Buffer: []byte{4, 240, 159, 140, 190, 14, 104, 195, 184, 115, 116, 101, 110, 32, 101, 114, 32, 102, 105, 110}})

		state.Rewind()
		value, err := Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "🌾")

		value, err = Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "høsten er fin")

		_, err = Encoder.Decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
