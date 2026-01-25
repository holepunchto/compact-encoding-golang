package compactencoding

type Bool struct{}

func (i *Bool) Preencode(state *State, _value bool) {
	state.End++
}

func (i *Bool) Encode(state *State, value bool) error {
	if value {
		state.Buffer[0] = 1
	} else {
		state.Buffer[0] = 0
	}

	state.Start = 1

	return nil
}

func (i *Bool) Decode(state *State) (bool, error) {
	if state.Start >= state.End {
		return false, &EncodingErrorOutOfBounds{}
	}

	return state.Buffer[0] == 1, nil
}

func NewBool() *Bool {
	return &Bool{}
}
