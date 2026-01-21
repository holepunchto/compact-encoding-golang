package compactencoding

type Uint8 struct{}

func (u *Uint8) preencode(state *State) {
	state.end += 1
}

func (u *Uint8) encode(state *State, value uint8) error {
	if state.start >= state.end {
		return &EncodingErrorOutOfBounds{}
	}

	state.buffer = make([]byte, state.end)
	state.buffer[state.start] = value
	state.start += 1

	return nil
}

func (u *Uint8) decode(state *State) (uint8, error) {
	if state.start >= state.end {
		return 0, &EncodingErrorOutOfBounds{}
	}

	value := state.buffer[state.start]
	state.start += 1

	return value, nil
}

func NewUint8() *Uint8 {
	return &Uint8{}
}

type Uint16 struct{}

func (u *Uint16) preencode(state *State) {
	state.end += 2
}

func (u *Uint16) encode(state *State, value uint16) error {
	if state.start >= state.end {
		return &EncodingErrorOutOfBounds{}
	}

	state.buffer[state.start] = byte(value & 0xff)
	state.buffer[state.start+1] = byte((value >> 8) & 0xff)
	state.start += 2

	return nil
}

func (u *Uint16) decode(state *State) (uint16, error) {
	if state.start >= state.end {
		return 0, &EncodingErrorOutOfBounds{}
	}

	value := uint16(state.buffer[state.start]) | uint16(state.buffer[state.start+1])<<8
	state.start += 2

	return value, nil
}

func NewUint16() *Uint16 {
	return &Uint16{}
}

type Uint32 struct{}

func (u *Uint32) preencode(state *State) {
	state.end += 4
}

func (u *Uint32) encode(state *State, value uint32) error {
	if state.start+4 > state.end {
		return &EncodingErrorOutOfBounds{}
	}
	state.buffer[state.start] = byte(value)
	state.buffer[state.start+1] = byte(value >> 8)
	state.buffer[state.start+2] = byte(value >> 16)
	state.buffer[state.start+3] = byte(value >> 24)
	state.start += 4
	return nil
}

func (u *Uint32) decode(state *State) (uint32, error) {
	if state.start+4 > state.end {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value := uint32(state.buffer[state.start]) |
		uint32(state.buffer[state.start+1])<<8 |
		uint32(state.buffer[state.start+2])<<16 |
		uint32(state.buffer[state.start+3])<<24
	state.start += 4
	return value, nil
}

func NewUint32() *Uint32 {
	return &Uint32{}
}

type Uint64 struct{}

func (u *Uint64) preencode(state *State) {
	state.end += 8
}

func (u *Uint64) encode(state *State, value uint64) error {
	if state.start+8 > state.end {
		return &EncodingErrorOutOfBounds{}
	}
	state.buffer[state.start] = byte(value)
	state.buffer[state.start+1] = byte(value >> 8)
	state.buffer[state.start+2] = byte(value >> 16)
	state.buffer[state.start+3] = byte(value >> 24)
	state.buffer[state.start+4] = byte(value >> 32)
	state.buffer[state.start+5] = byte(value >> 40)
	state.buffer[state.start+6] = byte(value >> 48)
	state.buffer[state.start+7] = byte(value >> 56)
	state.start += 8
	return nil
}

func (u *Uint64) decode(state *State) (uint64, error) {
	if state.start+8 > state.end {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value := uint64(state.buffer[state.start]) |
		uint64(state.buffer[state.start+1])<<8 |
		uint64(state.buffer[state.start+2])<<16 |
		uint64(state.buffer[state.start+3])<<24 |
		uint64(state.buffer[state.start+4])<<32 |
		uint64(state.buffer[state.start+5])<<40 |
		uint64(state.buffer[state.start+6])<<48 |
		uint64(state.buffer[state.start+7])<<56
	state.start += 8
	return value, nil
}

func NewUint64() *Uint64 {
	return &Uint64{}
}

type Uint struct{}

func (u *Uint) preencode(state *State, value uint) {
	if value <= 0xfc {
		state.end += 1
	} else if value <= 0xffff {
		state.end += 3
	} else if value <= 0xffffffff {
		state.end += 5
	} else {
		state.end += 9
	}
}

func (u *Uint) encode(state *State, value uint) error {
	if value <= 0xfc {
		return NewUint8().encode(state, uint8(value))
	}
	if state.start >= state.end {
		return &EncodingErrorOutOfBounds{}
	}
	if value <= 0xffff {
		state.buffer[state.start] = 0xfd
		state.start += 1
		return NewUint16().encode(state, uint16(value))
	}
	if value <= 0xffffffff {
		state.buffer[state.start] = 0xfe
		state.start += 1
		return NewUint32().encode(state, uint32(value))
	}
	state.buffer[state.start] = 0xff
	state.start += 1
	return NewUint64().encode(state, uint64(value))
}

func (u *Uint) decode(state *State) (uint, error) {
	if state.start >= state.end {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value, err := NewUint8().decode(state)
	if err != nil {
		return 0, err
	}
	if value <= 0xfc {
		return uint(value), nil
	}
	if value == 0xfd {
		v, err := NewUint16().decode(state)
		return uint(v), err
	}
	if value == 0xfe {
		v, err := NewUint32().decode(state)
		return uint(v), err
	}

	v, err := NewUint64().decode(state)
	return uint(v), err
}

func NewUint() *Uint {
	return &Uint{}
}
