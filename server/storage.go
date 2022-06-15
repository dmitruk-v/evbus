package server

type storage struct {
	idx  int
	msgs []*message
}

func NewStorage(size int) *storage {
	return &storage{
		idx:  0,
		msgs: make([]*message, size),
	}
}

func (st *storage) add(msg *message) {
	if st.idx == len(st.msgs) {
		st.idx = 0
	}
	st.msgs[st.idx] = msg
	st.idx++
}

func (st *storage) getAll() []*message {
	return st.msgs
}
