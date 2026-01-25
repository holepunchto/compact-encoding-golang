package compactencoding

type Int8 struct{}

func (i *Int8) Preencode(state *State, _value int8) {
	state.End += 1
}

func (i *Int8) Encode(state *State, value int8) error {
	Encoded := uint8((value << 1) ^ (value >> 7))
	return NewUint8().Encode(state, Encoded)
}

func (i *Int8) Decode(state *State) (int8, error) {
	v, err := NewUint8().Decode(state)
	if err != nil {
		return 0, err
	}
	return int8((v >> 1) ^ -(v & 1)), nil
}

func NewInt8() *Int8 {
	return &Int8{}
}

type Int16 struct{}

func (i *Int16) Preencode(state *State, _value int16) {
	state.End += 2
}

func (i *Int16) Encode(state *State, value int16) error {
	Encoded := uint16((value << 1) ^ (value >> 15))
	return NewUint16().Encode(state, Encoded)
}

func (i *Int16) Decode(state *State) (int16, error) {
	v, err := NewUint16().Decode(state)
	if err != nil {
		return 0, err
	}
	return int16((v >> 1) ^ -(v & 1)), nil
}

func NewInt16() *Int16 {
	return &Int16{}
}

type Int32 struct{}

func (i *Int32) Preencode(state *State, _value int32) {
	state.End += 4
}

func (i *Int32) Encode(state *State, value int32) error {
	Encoded := uint32((value << 1) ^ (value >> 31))
	return NewUint32().Encode(state, Encoded)
}

func (i *Int32) Decode(state *State) (int32, error) {
	v, err := NewUint32().Decode(state)
	if err != nil {
		return 0, err
	}
	return int32((v >> 1) ^ -(v & 1)), nil
}

func NewInt32() *Int32 {
	return &Int32{}
}

type Int64 struct{}

func (i *Int64) Preencode(state *State, _value int64) {
	state.End += 8
}

func (i *Int64) Encode(state *State, value int64) error {
	Encoded := uint64((value << 1) ^ (value >> 63))
	return NewUint64().Encode(state, Encoded)
}

func (i *Int64) Decode(state *State) (int64, error) {
	v, err := NewUint64().Decode(state)
	if err != nil {
		return 0, err
	}
	return int64((v >> 1) ^ -(v & 1)), nil
}

func NewInt64() *Int64 {
	return &Int64{}
}

type Int struct{}

func (i *Int) Preencode(state *State, value int) {
	Encoded := uint((value << 1) ^ (value >> 63))
	NewUint().Preencode(state, Encoded)
}

func (i *Int) Encode(state *State, value int) error {
	Encoded := uint((value << 1) ^ (value >> 63))
	return NewUint().Encode(state, Encoded)
}

func (i *Int) Decode(state *State) (int, error) {
	v, err := NewUint().Decode(state)
	if err != nil {
		return 0, err
	}
	return int((v >> 1) ^ -(v & 1)), nil
}

func NewInt() *Int {
	return &Int{}
}
