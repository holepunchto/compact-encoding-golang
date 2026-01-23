package compactencoding

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuffer(t *testing.T) {
	state := &State{}
	encoder := NewBuffer()

	// state.start = 0
	// t.alike(enc.buffer.decode(state), b4a.from('hi'))
	// t.alike(enc.buffer.decode(state), b4a.from('hello'))
	// t.is(enc.buffer.decode(state), null)
	// t.is(state.start, state.end)

	// t.exception(() => enc.buffer.decode(state))

	Convey("buffer can be encoded", t, func() {
		encoder.preencode(state, []byte("hi"))
		So(state, ShouldResemble, &State{start: 0, end: 3, buffer: nil})
		encoder.preencode(state, []byte("hello"))
		So(state, ShouldResemble, &State{start: 0, end: 9, buffer: nil})
		encoder.preencode(state, nil)
		So(state, ShouldResemble, &State{start: 0, end: 10, buffer: nil})

		state.Allocate()

		err := encoder.encode(state, []byte("hi"))
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 3, end: 10, buffer: []byte{2, 104, 105, 0, 0, 0, 0, 0, 0, 0}})

		err = encoder.encode(state, []byte("hello"))
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 9, end: 10, buffer: []byte{2, 104, 105, 5, 104, 101, 108, 108, 111, 0}})

		err = encoder.encode(state, nil)
		So(err, ShouldBeNil)
		So(state, ShouldResemble, &State{start: 10, end: 10, buffer: []byte{2, 104, 105, 5, 104, 101, 108, 108, 111, 0}})

		state.Rewind()

		value, err := encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, []byte("hi"))

		value, err = encoder.decode(state)
		So(err, ShouldBeNil)
		So(value, ShouldEqual, []byte("hello"))

		value, err = encoder.decode(state)
		So(err, ShouldBeNil)
		So(value == nil, ShouldBeTrue) // goconvey returns []byte(nil)

		_, err = encoder.decode(state)
		So(err.Error(), ShouldEqual, "EncodingError: Out of Bounds")
	})
}
