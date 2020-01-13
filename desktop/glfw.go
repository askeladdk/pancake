// +build darwin freebsd linux windows
// +build !js
// +build !android
// +build !ios

package desktop

import (
	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/graphics/opengl"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type window struct {
	window *glfw.Window
}

func (this *window) ShouldClose() bool {
	return this.window.ShouldClose()
}

func (this *window) Update() {
	mainthread.Call(func() {
		this.window.SwapBuffers()
		glfw.PollEvents()
	})
}

func (this *window) Size() (int, int) {
	return this.window.GetSize()
}

func newWindow(opt pancake.WindowOptions) (*window, error) {
	// glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if win, err := glfw.CreateWindow(opt.Width, opt.Height, opt.Title, nil, nil); err != nil {
		return nil, err
	} else {
		win.MakeContextCurrent()
		return &window{win}, nil
	}
}

func Run(opt pancake.WindowOptions, run func(pancake.Window)) {
	if err := glfw.Init(); err != nil {
		panic(err)
	} else if wnd, err := newWindow(opt); err != nil {
		panic(err)
	} else if err := opengl.Init(nil); err != nil {
		panic(err)
	} else {
		defer glfw.Terminate()
		mainthread.Run(func() {
			run(wnd)
		})
	}
}
