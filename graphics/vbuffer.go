package graphics

import (
	"errors"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var vbobinder = newBinder(func(id uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer(id))
})

type Buffer struct {
	format AttribFormat
	stride int
	count  int
	id     gl.Buffer
}

func (vbo *Buffer) Begin() {
	vbobinder.bind(uint32(vbo.id))
}

func (vbo *Buffer) End() {
	vbobinder.unbind()
}

func (vbo *Buffer) Len() int {
	return vbo.count
}

func (vbo *Buffer) AttribFormat() AttribFormat {
	return vbo.format
}

func (vbo *Buffer) SetData(i, j int, data interface{}) {
	if j > i && i >= 0 && j < vbo.count {
		gl.BufferSubData(gl.ARRAY_BUFFER, i*vbo.stride, (j-i)*vbo.stride, gl.Ptr(data))
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (vbo *Buffer) Draw(mode gl.Enum, i, j int) {
	if j > i && i >= 0 && j <= vbo.count {
		gl.DrawArrays(mode, i, j-i)
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (vbo *Buffer) delete() {
	gl.DeleteBuffer(vbo.id)
}

func NewBuffer(format AttribFormat, count int, data interface{}) *Buffer {
	buf := &Buffer{
		format: format,
		stride: format.stride(),
		count:  count,
		id:     gl.CreateBuffer(),
	}

	runtime.SetFinalizer(buf, (*Buffer).delete)

	buf.Begin()
	defer buf.End()

	gl.BufferData(gl.ARRAY_BUFFER, buf.stride*buf.count, gl.Ptr(data), gl.DYNAMIC_DRAW)
	return buf
}
