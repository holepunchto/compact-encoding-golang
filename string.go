package compactencoding

import (
	"strings"
	"unicode/utf8"
)

// String writes the length of a string in bytes followed by its UTF-8 bytes.
//
// As in JavaScript, a string is valid UTF-8 once encoded or decoded: Each
// maximal subpart of an ill-formed sequence becomes U+FFFD, as the Unicode
// Standard recommends in section 3.9.
var String Codec[string] = stringCodec{}

type stringCodec struct{}

func (stringCodec) Preencode(state *State, value string) {
	n := len(toValidUTF8(value))

	state.End += varintSize(uint64(n)) + n
}

func (stringCodec) Encode(state *State, value string) error {
	value = toValidUTF8(value)

	if err := encodeVarint(state, uint64(len(value))); err != nil {
		return err
	}

	b, err := state.take(len(value))
	if err != nil {
		return err
	}

	copy(b, value)

	return nil
}

func (stringCodec) Decode(state *State) (string, error) {
	n, err := decodeLength(state)
	if err != nil {
		return "", err
	}

	b, err := state.take(n)
	if err != nil {
		return "", err
	}

	return toValidUTF8(string(b)), nil
}

func toValidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	var b strings.Builder

	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])

		if r == utf8.RuneError && n == 1 {
			b.WriteRune(utf8.RuneError)

			i += maximalSubpart(s[i:])
		} else {
			b.WriteString(s[i : i+n])

			i += n
		}
	}

	return b.String()
}

// Returns the length of the longest prefix of s that starts a well-formed
// sequence, or 1 if none does, given that s starts an ill-formed one. The
// ranges are those of table 3-7 of the Unicode Standard.
func maximalSubpart(s string) int {
	var lo, hi byte = 0x80, 0xbf

	n := 0

	switch c := s[0]; {
	case c >= 0xc2 && c <= 0xdf:
		n = 2
	case c == 0xe0:
		lo, n = 0xa0, 3
	case c == 0xed:
		hi, n = 0x9f, 3
	case c >= 0xe1 && c <= 0xef:
		n = 3
	case c == 0xf0:
		lo, n = 0x90, 4
	case c == 0xf4:
		hi, n = 0x8f, 4
	case c >= 0xf1 && c <= 0xf3:
		n = 4
	default:
		return 1
	}

	i := 1

	for i < n && i < len(s) && s[i] >= lo && s[i] <= hi {
		i++

		lo, hi = 0x80, 0xbf
	}

	return i
}
