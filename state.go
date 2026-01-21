package compactencoding

type State struct {
	end    uint
	start  uint
	buffer []byte
}

func (state *State) Rewind() {
	state.start = 0
}

func (state *State) Allocate() {
	state.buffer = make([]byte, state.end)
}
