package graphics

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var vaobinder = newBinder(func(id uint32) {
	gl.BindVertexArray(gl.VertexArray(id))
})

type varray struct {
	vbo *Buffer
	id  gl.VertexArray
}

func (vao *varray) begin() {
	vaobinder.bind(uint32(vao.id))
}

func (vao *varray) end() {
	vaobinder.unbind()
}

func (vao *varray) setData(i, j int, data interface{}) {
	vao.vbo.Begin()
	defer vao.vbo.End()
	vao.vbo.SetData(i, j, data)
}

func (vao *varray) slice(i, j int) *VertexSlice {
	if j > i && i >= 0 && j <= vao.vbo.Len() {
		return &VertexSlice{
			vao: vao,
			i:   i,
			j:   j,
		}
	} else {
		panic(fmt.Errorf("range out of bounds"))
	}
}

func (vao *varray) delete() {
	gl.DeleteVertexArray(vao.id)
}

type VertexSlice struct {
	vao  *varray
	i, j int
}

func (vbs *VertexSlice) Begin() {
	vbs.vao.begin()
}

func (vbs *VertexSlice) End() {
	vbs.vao.end()
}

func (vbs *VertexSlice) Buffer() *Buffer {
	return vbs.vao.vbo
}

func (vbs *VertexSlice) SetData(data interface{}) {
	vbs.vao.setData(vbs.i, vbs.j, data)
}

func (vbs *VertexSlice) Slice(i, j int) *VertexSlice {
	if j-i > vbs.Len() {
		panic(errors.New("range out of bounds"))
	} else {
		return vbs.vao.slice(vbs.i+i, vbs.i+j)
	}
}

func (vbs *VertexSlice) Draw(mode gl.Enum) {
	vbs.vao.vbo.Draw(mode, vbs.i, vbs.j)
}

func (vbs *VertexSlice) Len() int {
	return vbs.j - vbs.i
}

func vaoAssignAttribs(buf *Buffer, attrib gl.Attrib) gl.Attrib {
	buf.Begin()
	defer buf.End()

	stride := buf.AttribFormat().stride()
	offset := 0

	for _, attr := range buf.AttribFormat() {
		for i := 0; i < attr.repeat(); i++ {
			gl.VertexAttribPointer(
				attrib, attr.components(), gl.FLOAT, false, stride, offset)
			// gl.VertexAttribDivisor(attrib, divisor)
			gl.EnableVertexAttribArray(attrib)
			offset += attr.stride()
			attrib += 1
		}
	}

	return attrib
}

func NewVertexSlice(vbo *Buffer) *VertexSlice {
	vao := &varray{
		vbo: vbo,
		id:  gl.CreateVertexArray(),
	}

	runtime.SetFinalizer(vao, (*varray).delete)

	vao.begin()
	defer vao.end()

	vaoAssignAttribs(vao.vbo, 0)

	return vao.slice(0, vbo.Len())
}
