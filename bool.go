package compactencoding

// Bool writes a single byte of 1 or 0. It decodes any byte but 1 as false.
var Bool Codec[bool] = boolCodec{}

type boolCodec struct{}

func (boolCodec) Preencode(state *State, _ bool) {
	state.End++
}

func (boolCodec) Encode(state *State, value bool) error {
	var b uint8

	if value {
		b = 1
	}

	return Uint8.Encode(state, b)
}

func (boolCodec) Decode(state *State) (bool, error) {
	b, err := Uint8.Decode(state)

	return b == 1, err
}
