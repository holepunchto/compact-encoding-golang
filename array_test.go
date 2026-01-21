package compactencoding

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArray(t *testing.T) {
	state := &State{}
	encoder := NewArray(NewUint())

	Convey("uint8 array can be encoded", t, func() {
		encoder.preencode(state, []uint{1, 2, 3})
		state.Allocate()

		err := encoder.encode(state, []uint{1, 2, 3})

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{3, 1, 2, 3})
	})

	Convey("uint8 array can be decoded", t, func() {
		state.Rewind()
		fmt.Println(state)
		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, []uint{1, 2, 3})
	})
}
