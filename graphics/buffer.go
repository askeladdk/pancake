package graphics

import (
	"errors"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var vbobinder = newBinder(func(ref uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer(ref))
})

type buffer struct {
	stride int
	len    int
	ref    gl.Buffer
}

func (this *buffer) Begin() {
	vbobinder.bind(uint32(this.ref))
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
		gl.BufferSubData(gl.ARRAY_BUFFER, i*this.stride, (j-i)*this.stride, gl.Ptr(data))
		panicError()
	}
}

func (this *buffer) delete() {
	gl.DeleteBuffer(this.ref)
}

func NewBuffer(stride, len int, data interface{}) (Buffer, error) {
	buf := &buffer{
		ref:    gl.CreateBuffer(),
		stride: stride,
		len:    len,
	}

	runtime.SetFinalizer(buf, (*buffer).delete)

	buf.Begin()
	defer buf.End()

	gl.BufferData(gl.ARRAY_BUFFER, buf.stride*buf.len, gl.Ptr(data), gl.DYNAMIC_DRAW)
	return buf, checkError()
}
