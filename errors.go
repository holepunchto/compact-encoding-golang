package compactencoding

type EncodingErrorOutOfBounds struct{}

func (e *EncodingErrorOutOfBounds) Error() string {
	return "EncodingError: Out of Bounds"
}

func NewEncodingErrorOutOfBounds() *EncodingErrorOutOfBounds {
	return &EncodingErrorOutOfBounds{}
}

type EncodingErrorArrayTooBig struct{}

func (e *EncodingErrorArrayTooBig) Error() string {
	return "EncodingError: Array is too big"
}

func NewEncodingErrorArrayTooBig() *EncodingErrorArrayTooBig {
	return &EncodingErrorArrayTooBig{}
}
