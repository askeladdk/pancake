// +build darwin freebsd linux windows
// +build !js
// +build !android
// +build !ios

package pancake

import (
	"image"

	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/input"

	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type window struct {
	window            *glfw.Window
	charEventHandler  input.CharEventHandler
	keyEventHandler   input.KeyEventHandler
	mouseEventHandler input.MouseEventHandler
	mousePosition     image.Point
	cursorEntered     bool
}

func (wnd *window) ShouldClose() bool {
	return wnd.window.ShouldClose()
}

func (wnd *window) PollEvents() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})
}

func (wnd *window) SwapBuffers() {
	mainthread.Call(func() {
		wnd.window.SwapBuffers()
	})
}

func (wnd *window) Bounds() image.Rectangle {
	var w, h int
	mainthread.Call(func() {
		w, h = wnd.window.GetSize()
	})
	return image.Rectangle{image.Point{}, image.Point{w, h}}
}

func (wnd *window) SetTitle(title string) {
	mainthread.Call(func() {
		wnd.window.SetTitle(title)
	})
}

func (wnd *window) SetCharEventHandler(handler input.CharEventHandler) {
	mainthread.Call(func() {
		wnd.charEventHandler = handler
		if handler != nil {
			wnd.window.SetCharCallback(wnd.charCallback)
		} else {
			wnd.window.SetCharCallback(nil)
		}
	})
}

func (wnd *window) SetKeyEventHandler(handler input.KeyEventHandler) {
	mainthread.Call(func() {
		wnd.keyEventHandler = handler
		if handler != nil {
			wnd.window.SetKeyCallback(wnd.keyCallback)
		} else {
			wnd.window.SetKeyCallback(nil)
		}
	})
}

func (wnd *window) SetMouseEventHandler(handler input.MouseEventHandler) {
	mainthread.Call(func() {
		wnd.mouseEventHandler = handler
		if handler != nil {
			wnd.window.SetCursorPosCallback(wnd.cursorCallback)
			wnd.window.SetCursorEnterCallback(wnd.cursorEnterCallback)
			wnd.window.SetMouseButtonCallback(wnd.mouseCallback)
		} else {
			wnd.window.SetCursorPosCallback(nil)
			wnd.window.SetCursorEnterCallback(nil)
			wnd.window.SetMouseButtonCallback(nil)
		}
	})
}

func makeInputFlags(action glfw.Action, mod glfw.ModifierKey) input.Flags {
	var flags input.Flags
	if action == glfw.Press {
		flags |= input.Pressed
	} else if action == glfw.Release {
		flags |= input.Released
	} else if action == glfw.Repeat {
		flags |= input.Repeated
	}

	if mod&glfw.ModAlt != 0 {
		flags |= input.Alt
	}

	if mod&glfw.ModControl != 0 {
		flags |= input.Control
	}

	if mod&glfw.ModShift != 0 {
		flags |= input.Shift
	}

	if mod&glfw.ModSuper != 0 {
		flags |= input.Super
	}

	return flags
}

func (wnd *window) charCallback(_ *glfw.Window, char rune) {
	wnd.charEventHandler.CharEvent(char)
}

func (wnd *window) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	wnd.keyEventHandler.KeyEvent(input.KeyEvent{
		Key:      input.Key(key),
		Flags:    makeInputFlags(action, mod),
		Scancode: scancode,
	})
}

func (wnd *window) cursorEnterCallback(_ *glfw.Window, entered bool) {
	wnd.cursorEntered = entered
}

func (wnd *window) cursorCallback(_ *glfw.Window, x, y float64) {
	if wnd.cursorEntered {
		wnd.mousePosition = image.Point{int(x), int(y)}
		wnd.mouseEventHandler.MouseEvent(input.MouseEvent{
			Mouse:    input.MouseMove,
			Position: wnd.mousePosition,
		})
	}
}

func (wnd *window) mouseCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if wnd.cursorEntered {
		wnd.mouseEventHandler.MouseEvent(input.MouseEvent{
			Mouse:    input.Mouse(button),
			Flags:    makeInputFlags(action, mod),
			Position: wnd.mousePosition,
		})
	}
}

func newWindow(opt WindowOptions) (*window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if glfwwnd, err := glfw.CreateWindow(opt.Size.X, opt.Size.Y, opt.Title, nil, nil); err != nil {
		return nil, err
	} else {
		wnd := &window{
			window: glfwwnd,
		}
		glfwwnd.MakeContextCurrent()
		return wnd, nil
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
