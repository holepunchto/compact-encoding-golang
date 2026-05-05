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
	if state.Start+uint(length) > state.End {
		return nil, &EncodingErrorOutOfBounds{}
	}
	value := state.Buffer[state.Start : state.Start+uint(length)]
	state.Start += uint(length)

	return value, nil
}

func NewBuffer() *Buffer {
	return &Buffer{}
}
