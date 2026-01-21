package compactencoding

type Encoder[T any] interface {
	preencode(state *State, value T)
	encode(state *State, value T) error
	decode(state *State) (T, error)
}

type Array[T any] struct {
	elementEncoder Encoder[T]
}

func (a *Array[T]) preencode(state *State, value []T) {
	NewUint().preencode(state, uint(len(value)))
	for _, e := range value {
		a.elementEncoder.preencode(state, e)
	}
}

func (a *Array[T]) encode(state *State, value []T) error {
	err := NewUint().encode(state, uint(len(value)))
	if err != nil {
		return err
	}
	for _, e := range value {
		err = a.elementEncoder.encode(state, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Array[T]) decode(state *State) ([]T, error) {
	length, err := NewUint().decode(state)
	if err != nil {
		return nil, err
	}
	result := make([]T, length)
	for i := range length {
		result[i], err = a.elementEncoder.decode(state)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func NewArray[T any](elementEncoder Encoder[T]) *Array[T] {
	return &Array[T]{elementEncoder: elementEncoder}
}
