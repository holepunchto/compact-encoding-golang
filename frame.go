package compactencoding

// Frame returns a codec that writes the encoded length of a value followed by
// the value, encoded with codec. Decoding reads the value from within the
// frame and skips any bytes of the frame that it leaves.
func Frame[T any](codec Codec[T]) Codec[T] {
	return frameCodec[T]{codec}
}

type frameCodec[T any] struct {
	codec Codec[T]
}

func (c frameCodec[T]) size(value T) int {
	var state State

	c.codec.Preencode(&state, value)

	return state.End
}

func (c frameCodec[T]) Preencode(state *State, value T) {
	n := c.size(value)

	state.End += varintSize(uint64(n)) + n
}

func (c frameCodec[T]) Encode(state *State, value T) error {
	if err := encodeVarint(state, uint64(c.size(value))); err != nil {
		return err
	}

	return c.codec.Encode(state, value)
}

func (c frameCodec[T]) Decode(state *State) (T, error) {
	n, err := decodeLength(state)
	if err != nil {
		var zero T

		return zero, err
	}

	end := state.End

	state.End = state.Start + n

	value, err := c.codec.Decode(state)

	state.Start, state.End = state.End, end

	return value, err
}
