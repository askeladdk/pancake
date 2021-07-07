package pancake

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/opengl"
)

var vaobinder = newBinder(func(id uint32) {
	gl.BindVertexArray(gl.VertexArray(id))
})

type vertexArrayObject struct {
	vbo *VertexBuffer
	id  gl.VertexArray
}

func (vao *vertexArrayObject) begin() {
	vaobinder.bind(uint32(vao.id))
}

func (vao *vertexArrayObject) end() {
	vaobinder.unbind()
}

func (vao *vertexArrayObject) setData(i, j int, data interface{}) {
	vao.vbo.Begin()
	defer vao.vbo.End()
	vao.vbo.SetData(i, j, data)
}

func (vao *vertexArrayObject) slice(i, j int) *VertexArraySlice {
	if j > i && i >= 0 && j <= vao.vbo.Len() {
		return &VertexArraySlice{
			vao: vao,
			i:   i,
			j:   j,
		}
	}
	panic(fmt.Errorf("range out of bounds"))
}

func (vao *vertexArrayObject) delete() {
	gl.DeleteVertexArray(vao.id)
}

type VertexArraySlice struct {
	vao  *vertexArrayObject
	i, j int
}

func (vas *VertexArraySlice) Begin() {
	vas.vao.begin()
}

func (vas *VertexArraySlice) End() {
	vas.vao.end()
}

func (vas *VertexArraySlice) Buffer() *VertexBuffer {
	return vas.vao.vbo
}

func (vas *VertexArraySlice) SetData(data interface{}) {
	vas.vao.setData(vas.i, vas.j, data)
}

func (vas *VertexArraySlice) Slice(i, j int) *VertexArraySlice {
	if j-i > vas.Len() {
		panic(errors.New("range out of bounds"))
	}
	return vas.vao.slice(vas.i+i, vas.i+j)
}

func (vas *VertexArraySlice) Draw(mode gl.Enum) {
	vas.vao.vbo.Draw(mode, vas.i, vas.j)
}

func (vas *VertexArraySlice) Len() int {
	return vas.j - vas.i
}

func vaoAssignAttribs(vbo *VertexBuffer, attrib gl.Attrib) gl.Attrib {
	vbo.Begin()
	defer vbo.End()

	stride := vbo.AttribFormat().stride()
	offset := 0

	for _, attr := range vbo.AttribFormat() {
		for i := 0; i < attr.repeat(); i++ {
			gl.VertexAttribPointer(
				attrib, attr.components(), attr.xtype(), attr.normalised(), stride, offset)
			// gl.VertexAttribDivisor(attrib, divisor)
			gl.EnableVertexAttribArray(attrib)
			offset += attr.stride()
			attrib++
		}
	}

	return attrib
}

func NewVertexArraySlice(vbo *VertexBuffer) *VertexArraySlice {
	vao := &vertexArrayObject{
		vbo: vbo,
		id:  gl.CreateVertexArray(),
	}

	runtime.SetFinalizer(vao, (*vertexArrayObject).delete)

	vao.begin()
	defer vao.end()

	vaoAssignAttribs(vao.vbo, 0)

	return vao.slice(0, vbo.Len())
}
