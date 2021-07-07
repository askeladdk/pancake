package pancake

import (
	"errors"
	"runtime"

	gl "github.com/askeladdk/pancake/opengl"
)

var vbobinder = newBinder(func(id uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer(id))
})

type VertexBuffer struct {
	format AttribFormat
	stride int
	count  int
	id     gl.Buffer
}

func (vbo *VertexBuffer) Begin() {
	vbobinder.bind(uint32(vbo.id))
}

func (vbo *VertexBuffer) End() {
	vbobinder.unbind()
}

func (vbo *VertexBuffer) Len() int {
	return vbo.count
}

func (vbo *VertexBuffer) AttribFormat() AttribFormat {
	return vbo.format
}

func (vbo *VertexBuffer) SetData(i, j int, data interface{}) {
	if j > i && i >= 0 && j <= vbo.count {
		gl.BufferSubData(gl.ARRAY_BUFFER, i*vbo.stride, (j-i)*vbo.stride, gl.Ptr(data))
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (vbo *VertexBuffer) Draw(mode gl.Enum, i, j int) {
	if j > i && i >= 0 && j <= vbo.count {
		gl.DrawArrays(mode, i, j-i)
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (vbo *VertexBuffer) delete() {
	gl.DeleteBuffer(vbo.id)
}

func NewVertexBuffer(format AttribFormat, count int, data interface{}) *VertexBuffer {
	buf := &VertexBuffer{
		format: format,
		stride: format.stride(),
		count:  count,
		id:     gl.CreateBuffer(),
	}

	runtime.SetFinalizer(buf, (*VertexBuffer).delete)

	buf.Begin()
	defer buf.End()

	gl.BufferData(gl.ARRAY_BUFFER, buf.stride*buf.count, gl.Ptr(data), gl.DYNAMIC_DRAW)
	return buf
}
