package compactencoding

type Buffer struct{}

func (b *Buffer) Preencode(state *State, value []byte) {
	if value != nil {
		NewUint().Preencode(state, uint(len(value)))
		state.End += uint(len(value))
	} else {
		state.End += 1
	}
}

func (b *Buffer) Encode(state *State, value []byte) error {
	if value != nil {
		err := NewUint().Encode(state, uint(len(value)))
		if err != nil {
			return err
		}
		if state.Start+uint(len(value)) > state.End {
			return &EncodingErrorOutOfBounds{}
		}
		copy(state.Buffer[state.Start:], value)
		state.Start += uint(len(value))
	} else {
		if state.Start >= state.End {
			return &EncodingErrorOutOfBounds{}
		}
		state.Buffer[state.Start] = 0
		state.Start += 1
	}
	return nil
}

func (b *Buffer) Decode(state *State) ([]byte, error) {
	length, err := NewUint().Decode(state)
	if err != nil {
		return nil, err
	}
	if state.End-state.Start < length {
		return nil, &EncodingErrorOutOfBounds{}
	}
	value := state.Buffer[state.Start : state.Start+length]
	state.Start += length

	return value, nil
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

// OptionalBuffer is the pre-3.0 buffer codec kept in JS as optionalBuffer: a
// nil slice encodes as a single zero byte and a zero-length payload decodes to
// nil. An empty non-nil slice therefore comes back as nil after a round trip.
type OptionalBuffer struct{}

func (b *OptionalBuffer) Preencode(state *State, value []byte) {
	NewBuffer().Preencode(state, value)
}

func (b *OptionalBuffer) Encode(state *State, value []byte) error {
	return NewBuffer().Encode(state, value)
}

func (b *OptionalBuffer) Decode(state *State) ([]byte, error) {
	value, err := NewBuffer().Decode(state)
	if err != nil || len(value) == 0 {
		return nil, err
	}
	return value, nil
}

func NewOptionalBuffer() *OptionalBuffer {
	return &OptionalBuffer{}
}
