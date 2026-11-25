package compactencoding

// The byte slice codecs, which write the length followed by the bytes.
// Decoding returns a slice of the buffer being decoded rather than a copy.
//
// Buffer encodes nil as it does an empty slice, and decodes an empty slice as
// empty rather than nil. OptionalBuffer decodes an empty slice as nil, so nil
// is what comes back from both.
var (
	Buffer         Codec[[]byte] = bufferCodec{}
	OptionalBuffer Codec[[]byte] = optionalBufferCodec{}
)

type bufferCodec struct{}

func (bufferCodec) Preencode(state *State, value []byte) {
	state.End += varintSize(uint64(len(value))) + len(value)
}

func (bufferCodec) Encode(state *State, value []byte) error {
	if err := encodeVarint(state, uint64(len(value))); err != nil {
		return err
	}

	b, err := state.take(len(value))
	if err != nil {
		return err
	}

	copy(b, value)

	return nil
}

func (bufferCodec) Decode(state *State) ([]byte, error) {
	n, err := decodeLength(state)
	if err != nil {
		return nil, err
	}

	return state.take(n)
}

type optionalBufferCodec struct {
	bufferCodec
}

func (c optionalBufferCodec) Decode(state *State) ([]byte, error) {
	b, err := c.bufferCodec.Decode(state)
	if err != nil || len(b) == 0 {
		return nil, err
	}

	return b, nil
}
