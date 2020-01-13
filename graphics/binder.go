package graphics

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

func (this *binder) bind(ref uint32) {
	this.stack = append(this.stack, this.cur)
	if ref != this.cur {
		this.bindfn(ref)
		this.cur = ref
	}
}

func (this *binder) unbind() {
	prev := this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	if prev != this.cur {
		this.bindfn(prev)
		this.cur = prev
	}
}
