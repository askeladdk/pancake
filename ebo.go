package pancake

import (
	"errors"
	"runtime"
	"unsafe"

	gl "github.com/askeladdk/pancake/opengl"
)

var ebobinder = newBinder(func(id uint32) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer(id))
})

type IndexBuffer struct {
	count int
	xtype gl.Enum
	id    gl.Buffer
}

func (ebo *IndexBuffer) Begin() {
	ebobinder.bind(uint32(ebo.id))
}

func (ebo *IndexBuffer) End() {
	ebobinder.unbind()
}

func (ebo *IndexBuffer) Len() int {
	return ebo.count
}

func (ebo *IndexBuffer) Draw(mode gl.Enum, i, j int) {
	if j > i && i >= 0 && j <= ebo.count {
		gl.DrawElements(mode, j-i, ebo.xtype, i)
		return
	}
	panic(errors.New("range out of bounds"))
}

func (ebo *IndexBuffer) delete() {
	gl.DeleteBuffer(ebo.id)
}

func newIndexBuffer(count, stride int, xtype gl.Enum, ptr unsafe.Pointer) *IndexBuffer {
	ebo := &IndexBuffer{
		count: count,
		xtype: xtype,
		id:    gl.CreateBuffer(),
	}

	runtime.SetFinalizer(ebo, (*IndexBuffer).delete)

	ebo.Begin()
	defer ebo.End()

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, count*stride, ptr, gl.STATIC_DRAW)
	return ebo
}

func NewIndexBufferUint8(indices []uint8) *IndexBuffer {
	return newIndexBuffer(len(indices), 1, gl.UNSIGNED_BYTE, gl.Ptr(indices))
}

func NewIndexBufferUint16(indices []uint16) *IndexBuffer {
	return newIndexBuffer(len(indices), 2, gl.UNSIGNED_SHORT, gl.Ptr(indices))
}

func NewIndexBufferUint32(indices []uint32) *IndexBuffer {
	return newIndexBuffer(len(indices), 4, gl.UNSIGNED_INT, gl.Ptr(indices))
}
