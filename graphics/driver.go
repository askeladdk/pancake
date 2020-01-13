package graphics

import (
	"fmt"
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type driver struct {
	screen frame
}

var theDriver *driver = nil

func Get() Driver {
	if theDriver != nil {
		return theDriver
	} else if err := gl.Init(); err != nil {
		panic(err)
	} else {
		theDriver = &driver{}
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		return theDriver
	}
}

func (this *driver) Version() string {
	return fmt.Sprint("OpenGL version ", gl.GoStr(gl.GetString(gl.VERSION)))
}

func (this *driver) Flush() {
	gl.Flush()
}

func (this *driver) SetViewport(bounds image.Rectangle) {
	size := bounds.Size()
	gl.Viewport(int32(bounds.Min.X), int32(bounds.Min.Y), int32(size.X), int32(size.Y))
}

func (this *driver) Clear(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
}

func (this *driver) NewTexture(size image.Point, filter Filter, format ColorFormat, pixels []byte) (Texture, error) {
	return newTexture(size, filter, format, pixels)
}

func (this *driver) NewShader(vshader, fshader string) (Shader, error) {
	return newShader(vshader, fshader)
}

func (this *driver) NewFrame(size image.Point, filter Filter, depthStencil bool) (Frame, error) {
	return newFrame(size, filter, depthStencil)
}

func (this *driver) ScreenFrame() Frame {
	return &this.screen
}

func (this *driver) NewVertexSlice(format AttrFormat, len int, data interface{}) (VertexSlice, error) {
	return newVao(format, len, data)
}
