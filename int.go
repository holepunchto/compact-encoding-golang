package compactencoding

import "math"

// The signed integer codecs, which zigzag the value into an unsigned one,
// mapping 0, -1, 1, -2 to 0, 1, 2, 3, and write it as the unsigned codec of
// the same width does.
var (
	Int   Codec[int]   = intCodec{}
	Int8  Codec[int8]  = int8Codec{}
	Int16 Codec[int16] = int16Codec{}
	Int32 Codec[int32] = int32Codec{}
	Int64 Codec[int64] = int64Codec{}
)

func zigzag(v int64) uint64 {
	return uint64(v<<1) ^ uint64(v>>63)
}

func unzigzag(v uint64) int64 {
	return int64(v>>1) ^ -int64(v&1)
}

type int8Codec struct{}

func (int8Codec) Preencode(state *State, _ int8) {
	state.End++
}

func (int8Codec) Encode(state *State, value int8) error {
	return Uint8.Encode(state, uint8(zigzag(int64(value))))
}

func (int8Codec) Decode(state *State) (int8, error) {
	v, err := Uint8.Decode(state)

	return int8(unzigzag(uint64(v))), err
}

type int16Codec struct{}

func (int16Codec) Preencode(state *State, _ int16) {
	state.End += 2
}

func (int16Codec) Encode(state *State, value int16) error {
	return Uint16.Encode(state, uint16(zigzag(int64(value))))
}

func (int16Codec) Decode(state *State) (int16, error) {
	v, err := Uint16.Decode(state)

	return int16(unzigzag(uint64(v))), err
}

type int32Codec struct{}

func (int32Codec) Preencode(state *State, _ int32) {
	state.End += 4
}

func (int32Codec) Encode(state *State, value int32) error {
	return Uint32.Encode(state, uint32(zigzag(int64(value))))
}

func (int32Codec) Decode(state *State) (int32, error) {
	v, err := Uint32.Decode(state)

	return int32(unzigzag(uint64(v))), err
}

type int64Codec struct{}

func (int64Codec) Preencode(state *State, _ int64) {
	state.End += 8
}

func (int64Codec) Encode(state *State, value int64) error {
	return Uint64.Encode(state, zigzag(value))
}

func (int64Codec) Decode(state *State) (int64, error) {
	v, err := Uint64.Decode(state)

	return unzigzag(v), err
}

type intCodec struct{}

func (intCodec) Preencode(state *State, value int) {
	state.End += varintSize(zigzag(int64(value)))
}

func (intCodec) Encode(state *State, value int) error {
	return encodeVarint(state, zigzag(int64(value)))
}

func (intCodec) Decode(state *State) (int, error) {
	v, err := decodeVarint(state)
	if err != nil {
		return 0, err
	}

	n := unzigzag(v)

	if n < math.MinInt || n > math.MaxInt {
		return 0, ErrOverflow
	}

	return int(n), nil
}
