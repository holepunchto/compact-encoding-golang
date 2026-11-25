package compactencoding

import "math"

// The floating-point codecs, which write the IEEE 754 binary32 and binary64
// forms in little-endian order.
var (
	Float32 Codec[float32] = float32Codec{}
	Float64 Codec[float64] = float64Codec{}
)

type float32Codec struct{}

func (float32Codec) Preencode(state *State, _ float32) {
	state.End += 4
}

func (float32Codec) Encode(state *State, value float32) error {
	return Uint32.Encode(state, math.Float32bits(value))
}

func (float32Codec) Decode(state *State) (float32, error) {
	v, err := Uint32.Decode(state)

	return math.Float32frombits(v), err
}

type float64Codec struct{}

func (float64Codec) Preencode(state *State, _ float64) {
	state.End += 8
}

func (float64Codec) Encode(state *State, value float64) error {
	return Uint64.Encode(state, math.Float64bits(value))
}

func (float64Codec) Decode(state *State) (float64, error) {
	v, err := Uint64.Decode(state)

	return math.Float64frombits(v), err
}
