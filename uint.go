package compactencoding

import (
	"encoding/binary"
	"math"
)

// The unsigned integer codecs. Uint writes a variable number of bytes, as
// few as the value needs, while the others write a fixed number of bytes in
// little-endian order.
var (
	Uint   Codec[uint]   = uintCodec{}
	Uint8  Codec[uint8]  = uint8Codec{}
	Uint16 Codec[uint16] = uint16Codec{}
	Uint32 Codec[uint32] = uint32Codec{}
	Uint64 Codec[uint64] = uint64Codec{}
)

type uint8Codec struct{}

func (uint8Codec) Preencode(state *State, _ uint8) {
	state.End++
}

func (uint8Codec) Encode(state *State, value uint8) error {
	b, err := state.take(1)
	if err != nil {
		return err
	}

	b[0] = value

	return nil
}

func (uint8Codec) Decode(state *State) (uint8, error) {
	b, err := state.take(1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

type uint16Codec struct{}

func (uint16Codec) Preencode(state *State, _ uint16) {
	state.End += 2
}

func (uint16Codec) Encode(state *State, value uint16) error {
	b, err := state.take(2)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(b, value)

	return nil
}

func (uint16Codec) Decode(state *State) (uint16, error) {
	b, err := state.take(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(b), nil
}

type uint32Codec struct{}

func (uint32Codec) Preencode(state *State, _ uint32) {
	state.End += 4
}

func (uint32Codec) Encode(state *State, value uint32) error {
	b, err := state.take(4)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(b, value)

	return nil
}

func (uint32Codec) Decode(state *State) (uint32, error) {
	b, err := state.take(4)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

type uint64Codec struct{}

func (uint64Codec) Preencode(state *State, _ uint64) {
	state.End += 8
}

func (uint64Codec) Encode(state *State, value uint64) error {
	b, err := state.take(8)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint64(b, value)

	return nil
}

func (uint64Codec) Decode(state *State) (uint64, error) {
	b, err := state.take(8)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(b), nil
}

type uintCodec struct{}

func (uintCodec) Preencode(state *State, value uint) {
	state.End += varintSize(uint64(value))
}

func (uintCodec) Encode(state *State, value uint) error {
	return encodeVarint(state, uint64(value))
}

func (uintCodec) Decode(state *State) (uint, error) {
	v, err := decodeVarint(state)
	if err != nil {
		return 0, err
	}

	if v > math.MaxUint {
		return 0, ErrOverflow
	}

	return uint(v), nil
}

// A varint is a single byte below 0xfd, or a byte of 0xfd, 0xfe, or 0xff
// followed by the value in 2, 4, or 8 bytes.
func varintSize(v uint64) int {
	switch {
	case v <= 0xfc:
		return 1
	case v <= 0xffff:
		return 3
	case v <= 0xffffffff:
		return 5
	default:
		return 9
	}
}

func encodeVarint(state *State, v uint64) error {
	b, err := state.take(varintSize(v))
	if err != nil {
		return err
	}

	switch len(b) {
	case 1:
		b[0] = byte(v)
	case 3:
		b[0] = 0xfd
		binary.LittleEndian.PutUint16(b[1:], uint16(v))
	case 5:
		b[0] = 0xfe
		binary.LittleEndian.PutUint32(b[1:], uint32(v))
	default:
		b[0] = 0xff
		binary.LittleEndian.PutUint64(b[1:], v)
	}

	return nil
}

func decodeVarint(state *State) (uint64, error) {
	prefix, err := Uint8.Decode(state)
	if err != nil {
		return 0, err
	}

	switch prefix {
	case 0xfd:
		v, err := Uint16.Decode(state)

		return uint64(v), err
	case 0xfe:
		v, err := Uint32.Decode(state)

		return uint64(v), err
	case 0xff:
		return Uint64.Decode(state)
	default:
		return uint64(prefix), nil
	}
}

// Decodes a length prefix, failing if fewer bytes than that are left.
func decodeLength(state *State) (int, error) {
	n, err := decodeVarint(state)
	if err != nil {
		return 0, err
	}

	if n > uint64(state.End-state.Start) {
		return 0, ErrOutOfBounds
	}

	return int(n), nil
}
