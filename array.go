package compactencoding

type Encoder[T any] interface {
	Preencode(state *State, value T)
	Encode(state *State, value T) error
	Decode(state *State) (T, error)
}

type Array[T any] struct {
	elementEncoder Encoder[T]
}

func (a *Array[T]) Preencode(state *State, value []T) {
	NewUint().Preencode(state, uint(len(value)))
	for _, e := range value {
		a.elementEncoder.Preencode(state, e)
	}
}

func (a *Array[T]) Encode(state *State, value []T) error {
	err := NewUint().Encode(state, uint(len(value)))
	if err != nil {
		return err
	}
	for _, e := range value {
		err = a.elementEncoder.Encode(state, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Array[T]) Decode(state *State) ([]T, error) {
	length, err := NewUint().Decode(state)
	if err != nil {
		return nil, err
	}
	result := make([]T, length)
	for i := range length {
		result[i], err = a.elementEncoder.Decode(state)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func NewArray[T any](elementEncoder Encoder[T]) *Array[T] {
	return &Array[T]{elementEncoder: elementEncoder}
}
