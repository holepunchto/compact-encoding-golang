package compactencoding

type String struct{}

func (s *String) Preencode(state *State, value string) {
	NewUint().Preencode(state, uint(len(value)))
	state.End += uint(len(value))
}

func (s *String) Encode(state *State, value string) error {
	count := uint(len(value))
	err := NewUint().Encode(state, uint(count))
	if err != nil {
		return err
	}

	if state.Start+count > state.End {
		return &EncodingErrorOutOfBounds{}
	}

	copy(state.Buffer[state.Start:], value)
	state.Start += count

	return nil
}

func (s *String) Decode(state *State) (string, error) {
	count, err := NewUint().Decode(state)
	if err != nil {
		return "", err
	}
	if state.Start+uint(int(count)) > state.End {
		return "", &EncodingErrorOutOfBounds{}
	}
	value := string(state.Buffer[state.Start : state.Start+uint(int(count))])
	state.Start += uint(int(count))
	return value, nil
}

func NewString() *String {
	return &String{}
}
