package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	state := &State{}
	encoder := NewString()

	Convey("string can be encoded", t, func() {
		encoder.preencode(state, "🌾")
		So(state, ShouldResemble, &State{start: 0, end: 5, buffer: nil})
		encoder.preencode(state, "høsten er fin")
		So(state, ShouldResemble, &State{start: 0, end: 20, buffer: nil})

		state.Allocate()
		err := encoder.encode(state, "🌾")
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 5, end: 20, buffer: []byte{4, 240, 159, 140, 190, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
		err = encoder.encode(state, "høsten er fin")
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 20, end: 20, buffer: []byte{4, 240, 159, 140, 190, 14, 104, 195, 184, 115, 116, 101, 110, 32, 101, 114, 32, 102, 105, 110}})

		state.Rewind()
		value, err := encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "🌾")

		value, err = encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "høsten er fin")

		_, err = encoder.decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
