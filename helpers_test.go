package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelpers(t *testing.T) {
	Encoder := NewInt()

	Convey("int can be Encoded", t, func() {
		buf, err := Encode(Encoder, 42)
		So(err, ShouldBeNil)
		So(buf, ShouldResemble, []byte{84})

		value, err := Decode(Encoder, buf)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 42)
	})
}
