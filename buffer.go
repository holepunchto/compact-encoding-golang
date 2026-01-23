package compactencoding

type Buffer struct{}

func (b *Buffer) preencode(state *State, value []byte) {
	if value != nil {
		NewUint().preencode(state, uint(len(value)))
		state.end += uint(len(value))
	} else {
		state.end += 1
	}
}

func (b *Buffer) encode(state *State, value []byte) error {
	if value != nil {
		err := NewUint().encode(state, uint(len(value)))
		if err != nil {
			return err
		}
		if state.start+uint(len(value)) > state.end {
			return &EncodingErrorOutOfBounds{}
		}
		copy(state.buffer[state.start:], value)
		state.start += uint(len(value))
	} else {
		if state.start >= state.end {
			return &EncodingErrorOutOfBounds{}
		}
		state.buffer[state.start] = 0
		state.start += 1
	}
	return nil
}

func (b *Buffer) decode(state *State) ([]byte, error) {
	length, err := NewUint().decode(state)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}
	if state.start+uint(length) > state.end {
		return nil, &EncodingErrorOutOfBounds{}
	}
	value := state.buffer[state.start : state.start+uint(length)]
	state.start += uint(length)

	return value, nil
}

func NewBuffer() *Buffer {
	return &Buffer{}
}
