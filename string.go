package compactencoding

type String struct{}

func (s *String) preencode(state *State, value string) {
	NewUint().preencode(state, uint(len(value)))
	state.end += uint(len(value))
}

func (s *String) encode(state *State, value string) error {
	count := uint(len(value))
	err := NewUint().encode(state, uint(count))
	if err != nil {
		return err
	}

	if state.start+count > state.end {
		return &EncodingErrorOutOfBounds{}
	}

	copy(state.buffer[state.start:], value)
	state.start += count

	return nil
}

func (s *String) decode(state *State) (string, error) {
	count, err := NewUint().decode(state)
	if err != nil {
		return "", err
	}
	if state.start+uint(int(count)) > state.end {
		return "", &EncodingErrorOutOfBounds{}
	}
	value := string(state.buffer[state.start : state.start+uint(int(count))])
	state.start += uint(int(count))
	return value, nil
}

func NewString() *String {
	return &String{}
}
