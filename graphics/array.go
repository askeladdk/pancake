package graphics

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var vaobinder = newBinder(func(ref uint32) {
	gl.BindVertexArray(ref)
})

func vaoAttribs(buf *buffer, format AttrFormat, atridx, divisor uint32) uint32 {
	buf.Begin()
	defer buf.End()

	stride := format.stride()
	offset := 0

	for _, attr := range format {
		for i := int32(0); i < attr.repeat(); i++ {
			gl.VertexAttribPointer(
				atridx, attr.components(), gl.FLOAT, false, stride, gl.PtrOffset(offset))
			gl.VertexAttribDivisor(atridx, divisor)
			gl.EnableVertexAttribArray(atridx)
			offset += int(attr.stride())
			atridx += 1
		}
	}

	return atridx
}

type vao struct {
	format AttrFormat
	shared *buffer
	ref    uint32
}

func (this *vao) Begin() {
	vaobinder.bind(this.ref)
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
			i:   int32(i),
			j:   int32(j),
		}
	} else {
		panic(fmt.Errorf("invalid range"))
	}
}

func (this *vao) draw(i, j int32) {
	if j > i && i >= 0 && j <= int32(this.Len()) {
		gl.DrawArrays(gl.TRIANGLES, int32(i), int32(j-i))
	} else {
		panic(errors.New("range out of bounds"))
	}
}

func (this *vao) Draw() {
	this.draw(0, int32(this.Len()))
}

func (this *vao) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteVertexArrays(1, &this.ref)
	})
}

func newVao(format AttrFormat, len int, data interface{}) (*vao, error) {
	if format == nil || len <= 0 {
		return nil, errors.New("invalid arguments")
	} else if shared, err := newBuffer(uint32(format.stride()), uint32(len), data); err != nil {
		return nil, err
	} else {
		v := &vao{
			format: format,
			shared: shared,
		}
		gl.GenVertexArrays(1, &v.ref)
		runtime.SetFinalizer(v, (*vao).delete)

		v.Begin()
		defer v.End()
		vaoAttribs(v.shared, format, 0, 0)
		return v, nil
	}
}

type vertexSlice struct {
	vao  *vao
	i, j int32
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
