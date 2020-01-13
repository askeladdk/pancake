package graphics

import (
	"errors"
	"runtime"

	"github.com/faiface/mainthread"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var vbobinder = newBinder(func(ref uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, ref)
})

type buffer struct {
	stride uint32
	len    uint32
	ref    uint32
}

func (this *buffer) Begin() {
	vbobinder.bind(this.ref)
}

func (this *buffer) End() {
	vbobinder.unbind()
}

func (this *buffer) Len() int {
	return int(this.len)
}

func (this *buffer) SetData(i, j int, data interface{}) {
	if j <= i || i < 0 || j > this.Len() {
		panic(errors.New("range out of bounds"))
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, i*int(this.stride), (j-i)*int(this.stride), gl.Ptr(data))
		panicError()
	}
}

func (this *buffer) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteBuffers(1, &this.ref)
	})
}

func newBuffer(stride, len uint32, data interface{}) (*buffer, error) {
	buf := &buffer{
		stride: stride,
		len:    len,
	}

	gl.GenBuffers(1, &buf.ref)
	runtime.SetFinalizer(buf, (*buffer).delete)

	buf.Begin()
	defer buf.End()

	gl.BufferData(gl.ARRAY_BUFFER, int(buf.stride*buf.len), gl.Ptr(data), gl.DYNAMIC_DRAW)
	return buf, checkError()
}
