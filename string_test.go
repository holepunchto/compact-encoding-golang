package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	state := &State{}
	encoder := NewString()

	Convey("string can be encoded", t, func() {
		encoder.preencode(state, "hello world!")
		state.Allocate()

		err := encoder.encode(state, "hello world!")

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{12, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33})
	})

	Convey("string can be decoded", t, func() {
		state.Rewind()
		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, "hello world!")
	})
}
