package graphics

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

type ivarray struct {
	vbo *Buffer
	ebo *IndexBuffer
	id  gl.VertexArray
}

func (vao *ivarray) begin() {
	vaobinder.bind(uint32(vao.id))
}

func (vao *ivarray) end() {
	vaobinder.unbind()
}

func (vao *ivarray) slice(i, j int) *IndexedVertexSlice {
	if j > i && i >= 0 && j <= vao.ebo.Len() {
		return &IndexedVertexSlice{
			vao: vao,
			i:   i,
			j:   j,
		}
	} else {
		panic(fmt.Errorf("range out of bounds"))
	}
}

func (vao *ivarray) delete() {
	gl.DeleteVertexArray(vao.id)
}

type IndexedVertexSlice struct {
	vao  *ivarray
	i, j int
}

func (vbs *IndexedVertexSlice) Begin() {
	vbs.vao.begin()
}

func (vbs *IndexedVertexSlice) End() {
	vbs.vao.end()
}

func (vbs *IndexedVertexSlice) Buffer() *Buffer {
	return vbs.vao.vbo
}

func (vbs *IndexedVertexSlice) Len() int {
	return vbs.j - vbs.i
}

func (vbs *IndexedVertexSlice) Draw(mode gl.Enum) {
	vbs.vao.ebo.Draw(mode, vbs.i, vbs.j)
}

func (vbs *IndexedVertexSlice) Slice(i, j int) *IndexedVertexSlice {
	if j-i > vbs.Len() {
		panic(errors.New("range out of bounds"))
	} else {
		return vbs.vao.slice(vbs.i+i, vbs.i+j)
	}
}

func NewIndexedVertexSlice(ebo *IndexBuffer, vbo *Buffer) *IndexedVertexSlice {
	vao := &ivarray{
		vbo: vbo,
		ebo: ebo,
		id:  gl.CreateVertexArray(),
	}

	runtime.SetFinalizer(vao, (*ivarray).delete)

	vao.begin()
	defer vao.end()

	// modifies the vao but not the global state
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.id)

	vaoAssignAttribs(vao.vbo, 0)

	return vao.slice(0, ebo.Len())
}
