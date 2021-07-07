package pancake

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/opengl"
)

type ivarray struct {
	vbo *VertexBuffer
	ebo *IndexBuffer
	id  gl.VertexArray
}

func (vao *ivarray) begin() {
	vaobinder.bind(uint32(vao.id))
}

func (vao *ivarray) end() {
	vaobinder.unbind()
}

func (vao *ivarray) slice(i, j int) *IndexedVertexArraySlice {
	if j > i && i >= 0 && j <= vao.ebo.Len() {
		return &IndexedVertexArraySlice{
			vao: vao,
			i:   i,
			j:   j,
		}
	}
	panic(fmt.Errorf("range out of bounds"))
}

func (vao *ivarray) delete() {
	gl.DeleteVertexArray(vao.id)
}

type IndexedVertexArraySlice struct {
	vao  *ivarray
	i, j int
}

func (ivas *IndexedVertexArraySlice) Begin() {
	ivas.vao.begin()
}

func (ivas *IndexedVertexArraySlice) End() {
	ivas.vao.end()
}

func (ivas *IndexedVertexArraySlice) VertexBuffer() *VertexBuffer {
	return ivas.vao.vbo
}

func (ivas *IndexedVertexArraySlice) Len() int {
	return ivas.j - ivas.i
}

func (ivas *IndexedVertexArraySlice) Draw(mode gl.Enum) {
	ivas.vao.ebo.Draw(mode, ivas.i, ivas.j)
}

func (ivas *IndexedVertexArraySlice) Slice(i, j int) *IndexedVertexArraySlice {
	if j-i > ivas.Len() {
		panic(errors.New("range out of bounds"))
	}
	return ivas.vao.slice(ivas.i+i, ivas.i+j)
}

func NewIndexedVertexArraySlice(ebo *IndexBuffer, vbo *VertexBuffer) *IndexedVertexArraySlice {
	iva := &ivarray{
		vbo: vbo,
		ebo: ebo,
		id:  gl.CreateVertexArray(),
	}

	runtime.SetFinalizer(iva, (*ivarray).delete)

	iva.begin()
	defer iva.end()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.id)

	vaoAssignAttribs(iva.vbo, 0)

	return iva.slice(0, ebo.Len())
}
