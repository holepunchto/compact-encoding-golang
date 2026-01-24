package compactencoding

type State struct {
	End    uint
	Start  uint
	Buffer []byte
}

func (state *State) Rewind() {
	state.Start = 0
}

func (state *State) Allocate() {
	state.Buffer = make([]byte, state.End)
}

func NewState(buf []byte) *State {
	return &State{End: uint(len(buf)), Buffer: buf}
}
