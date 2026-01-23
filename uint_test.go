package compactencoding

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUint8(t *testing.T) {
	state := &State{}
	encoder := NewUint8()

	Convey("uint8 can be encoded", t, func() {
		encoder.preencode(state)
		state.Allocate()

		err := encoder.encode(state, 12)

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{12})
	})

	Convey("uint8 can be decoded", t, func() {
		state.Rewind()
		fmt.Println(state)
		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, 12)
	})
}

func TestUint16(t *testing.T) {
	state := &State{}
	encoder := NewUint16()

	Convey("uint16 can be encoded", t, func() {
		encoder.preencode(state)
		state.Allocate()

		err := encoder.encode(state, 12|(34<<8))

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{12, 34})
	})

	Convey("uint16 can be decoded", t, func() {
		state.Rewind()

		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, 12|(34<<8))
	})
}

func TestUint32(t *testing.T) {
	state := &State{}
	encoder := NewUint32()

	Convey("uint32 can be encoded", t, func() {
		encoder.preencode(state)
		state.Allocate()

		err := encoder.encode(state, 12|(34<<8)|(56<<16)|(78<<24))

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{12, 34, 56, 78})
	})

	Convey("uint32 can be decoded", t, func() {
		state.Rewind()

		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, 12|(34<<8)|(56<<16)|(78<<24))
	})
}

func TestUint64(t *testing.T) {
	state := &State{}
	encoder := NewUint64()

	Convey("uint8 can be encoded", t, func() {
		encoder.preencode(state)
		state.Allocate()

		err := encoder.encode(state, 12|(34<<8)|(56<<16)|(78<<24)|(90<<32)|(12<<40)|(34<<48)|(56<<56))

		So(err, ShouldBeNil)
		So(state.buffer, ShouldEqual, []byte{12, 34, 56, 78, 90, 12, 34, 56})
	})

	Convey("uint32 can be decoded", t, func() {
		state.Rewind()

		value, err := encoder.decode(state)

		So(err, ShouldBeNil)
		So(value, ShouldEqual, 12|(34<<8)|(56<<16)|(78<<24)|(90<<32)|(12<<40)|(34<<48)|(56<<56))
	})
}

func TestUint(t *testing.T) {
	tests := []struct {
		Value    uint
		Encoding []byte
	}{
		{Value: 12, Encoding: []byte{12}},
		{Value: 12 | (34 << 8), Encoding: []byte{253, 12, 34}},
		{Value: 12 | (34 << 8) | (56 << 16) | (78 << 24), Encoding: []byte{254, 12, 34, 56, 78}},
		{Value: 12 | (34 << 8) | (56 << 16) | (78 << 24) | (90 << 32) | (12 << 40) | (34 << 48) | (56 << 56), Encoding: []byte{255, 12, 34, 56, 78, 90, 12, 34, 56}},
	}

	for _, v := range tests {
		Convey(fmt.Sprintf("uint with %v", v), t, func() {
			state := &State{}
			encoder := NewUint()

			encoder.preencode(state, v.Value)
			state.Allocate()
			err := encoder.encode(state, v.Value)

			So(err, ShouldBeNil)
			So(state.buffer, ShouldEqual, v.Encoding)

			state.Rewind()
			value, err := encoder.decode(state)
			So(err, ShouldBeNil)
			So(value, ShouldEqual, v.Value)
		})
	}
}
