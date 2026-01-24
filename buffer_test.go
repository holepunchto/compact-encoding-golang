package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuffer(t *testing.T) {
	state := &State{}
	Encoder := NewBuffer()

	// state.Start = 0
	// t.alike(enc.buffer.Decode(state), b4a.from('hi'))
	// t.alike(enc.buffer.Decode(state), b4a.from('hello'))
	// t.is(enc.buffer.Decode(state), null)
	// t.is(state.Start, state.End)

	// t.exception(() => enc.buffer.Decode(state))

	Convey("buffer can be Encoded", t, func() {
		Encoder.Preencode(state, []byte("hi"))
		So(state, ShouldResemble, &State{Start: 0, End: 3, Buffer: nil})
		Encoder.Preencode(state, []byte("hello"))
		So(state, ShouldResemble, &State{Start: 0, End: 9, Buffer: nil})
		Encoder.Preencode(state, nil)
		So(state, ShouldResemble, &State{Start: 0, End: 10, Buffer: nil})

		state.Allocate()

		err := Encoder.Encode(state, []byte("hi"))
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 3, End: 10, Buffer: []byte{2, 104, 105, 0, 0, 0, 0, 0, 0, 0}})

		err = Encoder.Encode(state, []byte("hello"))
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 9, End: 10, Buffer: []byte{2, 104, 105, 5, 104, 101, 108, 108, 111, 0}})

		err = Encoder.Encode(state, nil)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{Start: 10, End: 10, Buffer: []byte{2, 104, 105, 5, 104, 101, 108, 108, 111, 0}})

		state.Rewind()

		value, err := Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, []byte("hi"))

		value, err = Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, []byte("hello"))

		value, err = Encoder.Decode(state)
		So(err, ShouldBeNil)
		So(value == nil, ShouldBeTrue) // goconvey returns []byte(nil)

		_, err = Encoder.Decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
