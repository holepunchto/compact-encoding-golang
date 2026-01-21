package compactencoding

type EncodingErrorOutOfBounds struct{}

func (e *EncodingErrorOutOfBounds) Error() string {
	return "EncodingError: Out of Bounds"
}
