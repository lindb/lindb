package grouping

type Buffer struct {
	data []uint32

	w, r int
}

func (b *Buffer) Write(v uint32) {
	b.data = append(b.data, v)
	b.w++
}

func (b *Buffer) Read() uint32 {
	v := b.data[b.r]
	b.r++
	return v
}

func (b *Buffer) Reset() {
	b.data = b.data[0:0]
	b.w = 0
	b.r = 0
}
