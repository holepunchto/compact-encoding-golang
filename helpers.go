package compactencoding

func Decode[T any](encoder Encoder[T], data []byte) (T, error) {
	state := &State{}
	state.End = uint(len(data))
	state.Buffer = data

	return encoder.Decode(state)
}

func Encode[T any](encoder Encoder[T], value T) ([]byte, error) {
	state := &State{}
	encoder.Preencode(state, value)
	state.Allocate()
	if err := encoder.Encode(state, value); err != nil {
		return nil, err
	}

	return state.Buffer, nil
}
