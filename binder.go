package pancake

type binder struct {
	bindfn func(uint32)
	stack  []uint32
	cur    uint32
}

func newBinder(fn func(uint32)) *binder {
	return &binder{
		bindfn: fn,
	}
}

func (b *binder) bind(ref uint32) {
	b.stack = append(b.stack, b.cur)
	if ref != b.cur {
		b.bindfn(ref)
		b.cur = ref
	}
}

func (b *binder) unbind() {
	prev := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	if prev != b.cur {
		b.bindfn(prev)
		b.cur = prev
	}
}
