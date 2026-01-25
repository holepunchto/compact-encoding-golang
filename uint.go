package compactencoding

type Uint8 struct{}

func (u *Uint8) Preencode(state *State, _value uint8) {
	state.End += 1
}

func (u *Uint8) Encode(state *State, value uint8) error {
	if state.Start >= state.End {
		return &EncodingErrorOutOfBounds{}
	}

	state.Buffer[state.Start] = value
	state.Start += 1

	return nil
}

func (u *Uint8) Decode(state *State) (uint8, error) {
	if state.Start >= state.End {
		return 0, &EncodingErrorOutOfBounds{}
	}

	value := state.Buffer[state.Start]
	state.Start += 1

	return value, nil
}

func NewUint8() *Uint8 {
	return &Uint8{}
}

type Uint16 struct{}

func (u *Uint16) Preencode(state *State, _value uint16) {
	state.End += 2
}

func (u *Uint16) Encode(state *State, value uint16) error {
	if state.Start >= state.End {
		return &EncodingErrorOutOfBounds{}
	}

	state.Buffer[state.Start] = byte(value & 0xff)
	state.Buffer[state.Start+1] = byte((value >> 8) & 0xff)
	state.Start += 2

	return nil
}

func (u *Uint16) Decode(state *State) (uint16, error) {
	if state.Start >= state.End {
		return 0, &EncodingErrorOutOfBounds{}
	}

	value := uint16(state.Buffer[state.Start]) | uint16(state.Buffer[state.Start+1])<<8
	state.Start += 2

	return value, nil
}

func NewUint16() *Uint16 {
	return &Uint16{}
}

type Uint32 struct{}

func (u *Uint32) Preencode(state *State, _value uint32) {
	state.End += 4
}

func (u *Uint32) Encode(state *State, value uint32) error {
	if state.Start+4 > state.End {
		return &EncodingErrorOutOfBounds{}
	}
	state.Buffer[state.Start] = byte(value)
	state.Buffer[state.Start+1] = byte(value >> 8)
	state.Buffer[state.Start+2] = byte(value >> 16)
	state.Buffer[state.Start+3] = byte(value >> 24)
	state.Start += 4
	return nil
}

func (u *Uint32) Decode(state *State) (uint32, error) {
	if state.Start+4 > state.End {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value := uint32(state.Buffer[state.Start]) |
		uint32(state.Buffer[state.Start+1])<<8 |
		uint32(state.Buffer[state.Start+2])<<16 |
		uint32(state.Buffer[state.Start+3])<<24
	state.Start += 4
	return value, nil
}

func NewUint32() *Uint32 {
	return &Uint32{}
}

type Uint64 struct{}

func (u *Uint64) Preencode(state *State, _value uint64) {
	state.End += 8
}

func (u *Uint64) Encode(state *State, value uint64) error {
	if state.Start+8 > state.End {
		return &EncodingErrorOutOfBounds{}
	}
	state.Buffer[state.Start] = byte(value)
	state.Buffer[state.Start+1] = byte(value >> 8)
	state.Buffer[state.Start+2] = byte(value >> 16)
	state.Buffer[state.Start+3] = byte(value >> 24)
	state.Buffer[state.Start+4] = byte(value >> 32)
	state.Buffer[state.Start+5] = byte(value >> 40)
	state.Buffer[state.Start+6] = byte(value >> 48)
	state.Buffer[state.Start+7] = byte(value >> 56)
	state.Start += 8
	return nil
}

func (u *Uint64) Decode(state *State) (uint64, error) {
	if state.Start+8 > state.End {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value := uint64(state.Buffer[state.Start]) |
		uint64(state.Buffer[state.Start+1])<<8 |
		uint64(state.Buffer[state.Start+2])<<16 |
		uint64(state.Buffer[state.Start+3])<<24 |
		uint64(state.Buffer[state.Start+4])<<32 |
		uint64(state.Buffer[state.Start+5])<<40 |
		uint64(state.Buffer[state.Start+6])<<48 |
		uint64(state.Buffer[state.Start+7])<<56
	state.Start += 8
	return value, nil
}

func NewUint64() *Uint64 {
	return &Uint64{}
}

type Uint struct{}

func (u *Uint) Preencode(state *State, value uint) {
	if value <= 0xfc {
		state.End += 1
	} else if value <= 0xffff {
		state.End += 3
	} else if value <= 0xffffffff {
		state.End += 5
	} else {
		state.End += 9
	}
}

func (u *Uint) Encode(state *State, value uint) error {
	if value <= 0xfc {
		return NewUint8().Encode(state, uint8(value))
	}
	if state.Start >= state.End {
		return &EncodingErrorOutOfBounds{}
	}
	if value <= 0xffff {
		state.Buffer[state.Start] = 0xfd
		state.Start += 1
		return NewUint16().Encode(state, uint16(value))
	}
	if value <= 0xffffffff {
		state.Buffer[state.Start] = 0xfe
		state.Start += 1
		return NewUint32().Encode(state, uint32(value))
	}
	state.Buffer[state.Start] = 0xff
	state.Start += 1
	return NewUint64().Encode(state, uint64(value))
}

func (u *Uint) Decode(state *State) (uint, error) {
	if state.Start >= state.End {
		return 0, &EncodingErrorOutOfBounds{}
	}
	value, err := NewUint8().Decode(state)
	if err != nil {
		return 0, err
	}
	if value <= 0xfc {
		return uint(value), nil
	}
	if value == 0xfd {
		v, err := NewUint16().Decode(state)
		return uint(v), err
	}
	if value == 0xfe {
		v, err := NewUint32().Decode(state)
		return uint(v), err
	}

	v, err := NewUint64().Decode(state)
	return uint(v), err
}

func NewUint() *Uint {
	return &Uint{}
}
