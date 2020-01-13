package graphics

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var vaobinder = newBinder(func(ref uint32) {
	gl.BindVertexArray(gl.VertexArray(ref))
})

func vaoAttribs(buf *buffer, format AttrFormat, attrib gl.Attrib, divisor int) gl.Attrib {
	buf.Begin()
	defer buf.End()

	stride := format.stride()
	offset := 0

	for _, attr := range format {
		for i := 0; i < attr.repeat(); i++ {
			gl.VertexAttribPointer(
				attrib, attr.components(), gl.FLOAT, false, stride, offset)
			gl.VertexAttribDivisor(attrib, divisor)
			gl.EnableVertexAttribArray(attrib)
			offset += attr.stride()
			attrib += 1
		}
	}

	return attrib
}

type vao struct {
	format AttrFormat
	shared *buffer
	ref    gl.VertexArray
}

func (this *vao) Begin() {
	vaobinder.bind(uint32(this.ref))
}

func (this *vao) End() {
	vaobinder.unbind()
}

func (this *vao) Format() AttrFormat {
	return this.format
}

func (this *vao) Len() int {
	return this.shared.Len()
}

func (this *vao) setData(i, j int, data interface{}) {
	this.shared.Begin()
	defer this.shared.End()
	this.shared.SetData(i, j, data)
}

func (this *vao) SetData(data interface{}) {
	this.setData(0, this.Len(), data)
}

func (this *vao) Slice(i, j int) VertexSlice {
	if j > i && i >= 0 && j <= this.Len() {
		return &vertexSlice{
			vao: this,
			i:   i,
			j:   j,
		}
	} else {
		panic(fmt.Errorf("invalid range"))
	}
}

func (this *vao) draw(i, j int) {
	if j > i && i >= 0 && j <= this.Len() {
		gl.DrawArrays(gl.TRIANGLES, i, j-i)
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (this *vao) Draw() {
	this.draw(0, this.Len())
}

func (this *vao) delete() {
	gl.DeleteVertexArray(this.ref)
}

func NewVertexSlice(format AttrFormat, len int, data interface{}) (*vao, error) {
	if format == nil || len <= 0 {
		return nil, errors.New("invalid arguments")
	} else if shared, err := NewBuffer(format.stride(), len, data); err != nil {
		return nil, err
	} else {
		v := &vao{
			ref:    gl.CreateVertexArray(),
			format: format,
			shared: shared.(*buffer),
		}
		runtime.SetFinalizer(v, (*vao).delete)

		v.Begin()
		defer v.End()
		vaoAttribs(v.shared, format, 0, 0)
		return v, nil
	}
}

type vertexSlice struct {
	vao  *vao
	i, j int
}

func (this *vertexSlice) Begin() {
	this.vao.Begin()
}

func (this *vertexSlice) End() {
	this.vao.End()
}

func (this *vertexSlice) SetData(data interface{}) {
	this.vao.setData(int(this.i), int(this.j), data)
}

func (this *vertexSlice) Slice(i, j int) VertexSlice {
	if j-i > this.Len() {
		panic(errors.New("range out of bounds"))
	} else {
		return this.vao.Slice(int(this.i)+i, int(this.i)+j)
	}
}

func (this *vertexSlice) Draw() {
	this.vao.draw(this.i, this.j)
}

func (this *vertexSlice) Len() int {
	return int(this.j - this.i)
}

func (this *vertexSlice) Format() AttrFormat {
	return this.vao.Format()
}
