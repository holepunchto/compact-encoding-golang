package compactencoding

// MaxArrayLength is the most elements that an array decodes, so that a
// hostile length can't force a huge allocation.
const MaxArrayLength = 0x100000

// Array returns a codec that writes the length of a slice followed by its
// elements, each encoded with codec. Decoding an array of more than
// MaxArrayLength elements fails with ErrArrayTooBig.
func Array[T any](codec Codec[T]) Codec[[]T] {
	return arrayCodec[T]{codec}
}

type arrayCodec[T any] struct {
	codec Codec[T]
}

func (c arrayCodec[T]) Preencode(state *State, value []T) {
	state.End += varintSize(uint64(len(value)))

	for _, v := range value {
		c.codec.Preencode(state, v)
	}
}

func (c arrayCodec[T]) Encode(state *State, value []T) error {
	if err := encodeVarint(state, uint64(len(value))); err != nil {
		return err
	}

	for _, v := range value {
		if err := c.codec.Encode(state, v); err != nil {
			return err
		}
	}

	return nil
}

func (c arrayCodec[T]) Decode(state *State) ([]T, error) {
	n, err := decodeVarint(state)
	if err != nil {
		return nil, err
	}

	if n > MaxArrayLength {
		return nil, ErrArrayTooBig
	}

	value := make([]T, n)

	for i := range value {
		if value[i], err = c.codec.Decode(state); err != nil {
			return nil, err
		}
	}

	return value, nil
}
