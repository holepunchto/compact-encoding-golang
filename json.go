package compactencoding

import (
	"bytes"
	"encoding/json"
)

// JSON returns a codec that writes a value as a String of its JSON encoding,
// with encoding/json. HTML characters are left unescaped, as
// `JSON.stringify()` in JavaScript leaves them.
func JSON[T any]() Codec[T] {
	return jsonCodec[T]{}
}

type jsonCodec[T any] struct{}

func marshalJSON(value any) (string, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(value); err != nil {
		return "", err
	}

	return string(bytes.TrimSuffix(b.Bytes(), []byte("\n"))), nil
}

// A value that fails to marshal adds nothing, as Encode then fails.
func (jsonCodec[T]) Preencode(state *State, value T) {
	if s, err := marshalJSON(value); err == nil {
		String.Preencode(state, s)
	}
}

func (jsonCodec[T]) Encode(state *State, value T) error {
	s, err := marshalJSON(value)
	if err != nil {
		return err
	}

	return String.Encode(state, s)
}

func (jsonCodec[T]) Decode(state *State) (T, error) {
	var value, zero T

	s, err := String.Decode(state)
	if err != nil {
		return zero, err
	}

	if err := json.Unmarshal([]byte(s), &value); err != nil {
		return zero, err
	}

	return value, nil
}
