package compactencoding

type Int8 struct{}

func (i *Int8) preencode(state *State) {
	state.end += 1
}

func (i *Int8) encode(state *State, value int8) error {
	encoded := uint8((value << 1) ^ (value >> 7))
	return NewUint8().encode(state, encoded)
}

func (i *Int8) decode(state *State) (int8, error) {
	v, err := NewUint8().decode(state)
	if err != nil {
		return 0, err
	}
	return int8((v >> 1) ^ -(v & 1)), nil
}

func NewInt8() *Int8 {
	return &Int8{}
}

type Int16 struct{}

func (i *Int16) preencode(state *State) {
	state.end += 2
}

func (i *Int16) encode(state *State, value int16) error {
	encoded := uint16((value << 1) ^ (value >> 15))
	return NewUint16().encode(state, encoded)
}

func (i *Int16) decode(state *State) (int16, error) {
	v, err := NewUint16().decode(state)
	if err != nil {
		return 0, err
	}
	return int16((v >> 1) ^ -(v & 1)), nil
}

func NewInt16() *Int16 {
	return &Int16{}
}

type Int32 struct{}

func (i *Int32) preencode(state *State) {
	state.end += 4
}

func (i *Int32) encode(state *State, value int32) error {
	encoded := uint32((value << 1) ^ (value >> 31))
	return NewUint32().encode(state, encoded)
}

func (i *Int32) decode(state *State) (int32, error) {
	v, err := NewUint32().decode(state)
	if err != nil {
		return 0, err
	}
	return int32((v >> 1) ^ -(v & 1)), nil
}

func NewInt32() *Int32 {
	return &Int32{}
}

type Int64 struct{}

func (i *Int64) preencode(state *State) {
	state.end += 8
}

func (i *Int64) encode(state *State, value int64) error {
	encoded := uint64((value << 1) ^ (value >> 63))
	return NewUint64().encode(state, encoded)
}

func (i *Int64) decode(state *State) (int64, error) {
	v, err := NewUint64().decode(state)
	if err != nil {
		return 0, err
	}
	return int64((v >> 1) ^ -(v & 1)), nil
}

func NewInt64() *Int64 {
	return &Int64{}
}

type Int struct{}

func (i *Int) preencode(state *State, value int) {
	encoded := uint((value << 1) ^ (value >> 63))
	NewUint().preencode(state, encoded)
}

func (i *Int) encode(state *State, value int) error {
	encoded := uint((value << 1) ^ (value >> 63))
	return NewUint().encode(state, encoded)
}

func (i *Int) decode(state *State) (int, error) {
	v, err := NewUint().decode(state)
	if err != nil {
		return 0, err
	}
	return int((v >> 1) ^ -(v & 1)), nil
}

func NewInt() *Int {
	return &Int{}
}
