// Package compactencoding implements compact-encoding, a binary encoding that
// is wire compatible with the JavaScript compact-encoding library.
package compactencoding

import "errors"

var (
	// ErrOutOfBounds is returned when a value doesn't fit in the bytes that
	// are left, such as when decoding input that ends early.
	ErrOutOfBounds = errors.New("compactencoding: out of bounds")

	// ErrArrayTooBig is returned when decoding an array of more than
	// MaxArrayLength elements.
	ErrArrayTooBig = errors.New("compactencoding: array is too big")

	// ErrOverflow is returned when decoding an integer that the Go type
	// can't hold, which only happens for Uint and Int on 32-bit platforms.
	ErrOverflow = errors.New("compactencoding: integer overflows its type")
)

// A Codec encodes and decodes values of type T.
//
// Encoding takes two passes. Preencode adds the encoded size of a value to
// State.End, after which the buffer is allocated and Encode writes the value
// at State.Start. Decode reads a value at State.Start, reading no further
// than State.End, and leaves State.Start after the value.
type Codec[T any] interface {
	Preencode(state *State, value T)
	Encode(state *State, value T) error
	Decode(state *State) (T, error)
}

// A State is a position in a buffer that codecs encode to or decode from.
type State struct {
	Start  int
	End    int
	Buffer []byte
}

// Allocate allocates a buffer of End bytes, the size that preencoding summed.
func (s *State) Allocate() {
	s.Buffer = make([]byte, s.End)
}

func (s *State) take(n int) ([]byte, error) {
	if n < 0 || n > s.End-s.Start || s.Start+n > len(s.Buffer) {
		return nil, ErrOutOfBounds
	}

	b := s.Buffer[s.Start : s.Start+n]

	s.Start += n

	return b, nil
}

// Encode encodes value with codec.
func Encode[T any](codec Codec[T], value T) ([]byte, error) {
	var state State

	codec.Preencode(&state, value)
	state.Allocate()

	if err := codec.Encode(&state, value); err != nil {
		return nil, err
	}

	return state.Buffer, nil
}

// Decode decodes a value from b with codec, ignoring any bytes after it.
func Decode[T any](codec Codec[T], b []byte) (T, error) {
	return codec.Decode(&State{End: len(b), Buffer: b})
}
