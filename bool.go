package compactencoding

type Bool struct{}

func (i *Bool) Preencode(state *State, _value bool) {
	state.End += 1
}

func (i *Bool) Encode(state *State, value bool) error {
	if state.Start >= state.End {
		return &EncodingErrorOutOfBounds{}
	}

	if value {
		state.Buffer[state.Start] = 1
	}
	state.Start += 1

	return nil
}

func (i *Bool) Decode(state *State) (bool, error) {
	if state.Start >= state.End {
		return false, &EncodingErrorOutOfBounds{}
	}

	value := state.Buffer[state.Start]
	state.Start += 1

	return value == 1, nil
}

func NewBool() *Bool {
	return &Bool{}
}
