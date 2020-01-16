// +build darwin freebsd linux windows
// +build !js
// +build !android
// +build !ios

package pancake

import (
	"image"

	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type window struct {
	window *glfw.Window
}

func (this *window) ShouldClose() bool {
	return this.window.ShouldClose()
}

func (this *window) SwapBuffers() {
	mainthread.Call(func() {
		this.window.SwapBuffers()
		glfw.PollEvents()
	})
}

func (this *window) Bounds() image.Rectangle {
	var w, h int
	mainthread.Call(func() {
		w, h = this.window.GetSize()
	})
	return image.Rectangle{image.Point{}, image.Point{w, h}}
}

func (this *window) SetTitle(title string) {
	mainthread.Call(func() {
		this.window.SetTitle(title)
	})
}

func newWindow(opt WindowOptions) (*window, error) {
	// glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if win, err := glfw.CreateWindow(opt.Size.X, opt.Size.Y, opt.Title, nil, nil); err != nil {
		return nil, err
	} else {
		win.MakeContextCurrent()
		return &window{win}, nil
	}
}

func Run(opt WindowOptions, run func(Window)) {
	if err := glfw.Init(); err != nil {
		panic(err)
	} else if wnd, err := newWindow(opt); err != nil {
		panic(err)
	} else if err := gl.Init(nil); err != nil {
		panic(err)
	} else {
		defer glfw.Terminate()
		mainthread.Run(func() {
			run(wnd)
		})
	}
}
