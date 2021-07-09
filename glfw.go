package pancake

import (
	"context"
	"image"

	gl "github.com/askeladdk/pancake/opengl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.design/x/mainthread"
)

type glfwWindow struct {
	*glfw.Window
	inputEvents     []interface{}
	cursorEntered   bool
	windowScale     int
	resolutionScale int
	viewport        image.Rectangle
	resolution      image.Point
	mousePosition   image.Point
}

func newGlfwWindow(opt Options) (*glfwWindow, error) {
	var wnd glfwWindow
	var err error

	// create the window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	if wnd.Window, err = glfw.CreateWindow(opt.WindowSize.X, opt.WindowSize.Y, opt.Title, nil, nil); err != nil {
		return nil, err
	}

	// register callbacks
	wnd.SetCharCallback(wnd.charCallback)
	wnd.SetCursorPosCallback(wnd.cursorCallback)
	wnd.SetCursorEnterCallback(wnd.cursorEnterCallback)
	wnd.SetKeyCallback(wnd.keyCallback)
	wnd.SetMouseButtonCallback(wnd.mouseCallback)

	// set the window parameters
	w, h := wnd.GetFramebufferSize()
	wnd.resolution = opt.Resolution
	wnd.windowScale = w / opt.WindowSize.X
	wnd.viewport = logicalViewport(image.Pt(w, h), opt.Resolution)
	wnd.resolutionScale = wnd.viewport.Size().X / wnd.resolution.X

	wnd.MakeContextCurrent()
	return &wnd, nil
}

func (wnd *glfwWindow) Begin() {
	gl.Viewport(image.Rectangle{Max: image.Pt(wnd.GetFramebufferSize())})
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Viewport(wnd.viewport)
}

func (wnd *glfwWindow) End() {
	mainthread.Call(wnd.SwapBuffers)
}

func (wnd *glfwWindow) Scissor(r image.Rectangle) Scissor {
	r.Min = r.Min.Mul(wnd.resolutionScale)
	r.Max = r.Max.Mul(wnd.resolutionScale)
	r = r.Add(wnd.viewport.Min)
	return Scissor(r)
}

func (wnd *glfwWindow) SetTitle(title string) {
	mainthread.Call(func() {
		wnd.Window.SetTitle(title)
	})
}

func (wnd *glfwWindow) Bounds() image.Rectangle {
	return image.Rectangle{Max: wnd.resolution}
}

func (wnd *glfwWindow) InputEvents(ctx context.Context, ch chan<- interface{}) bool {
	if wnd.ShouldClose() {
		wnd.SetShouldClose(false)
		select {
		case <-ctx.Done():
			return false
		case ch <- QuitEvent{}:
		}
	}

	mainthread.Call(glfw.PollEvents)
	for _, inputEvent := range wnd.inputEvents {
		select {
		case <-ctx.Done():
			return false
		case ch <- inputEvent:
		}
	}
	wnd.inputEvents = wnd.inputEvents[:0]
	return true
}

func (wnd *glfwWindow) charCallback(_ *glfw.Window, char rune) {
	wnd.inputEvents = append(wnd.inputEvents, CharEvent{
		Char: char,
	})
}

func (wnd *glfwWindow) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	wnd.inputEvents = append(wnd.inputEvents, KeyEvent{
		Key:       Key(key),
		Modifiers: glfwModifiers(action, mod),
		Scancode:  scancode,
	})
}

func (wnd *glfwWindow) cursorEnterCallback(_ *glfw.Window, entered bool) {
	wnd.cursorEntered = entered
}

func (wnd *glfwWindow) cursorCallback(_ *glfw.Window, x, y float64) {
	mouse := image.Point{int(x), int(y)}.Mul(wnd.windowScale)

	if wnd.cursorEntered && mouse.In(wnd.viewport) {
		// Scale the mouse position from window to resolution coordinates.
		mouse = mouse.Sub(wnd.viewport.Min)
		vpsz := wnd.viewport.Size()
		wnd.mousePosition = image.Point{
			mouse.X * wnd.resolution.X / vpsz.X,
			mouse.Y * wnd.resolution.Y / vpsz.Y,
		}

		wnd.inputEvents = append(wnd.inputEvents, MouseMoveEvent{
			Position: wnd.mousePosition,
		})
	}
}

func (wnd *glfwWindow) mouseCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if wnd.cursorEntered {
		wnd.inputEvents = append(wnd.inputEvents, MouseEvent{
			Button:    MouseButton(button),
			Modifiers: glfwModifiers(action, mod),
			Position:  wnd.mousePosition,
		})
	}
}

func glfwModifiers(action glfw.Action, mod glfw.ModifierKey) (modifiers Modifiers) {
	if action == glfw.Press {
		modifiers |= ModPressed
	} else if action == glfw.Release {
		modifiers |= ModReleased
	} else if action == glfw.Repeat {
		modifiers |= ModRepeated
	}

	if mod&glfw.ModAlt != 0 {
		modifiers |= ModAlt
	}

	if mod&glfw.ModControl != 0 {
		modifiers |= ModControl
	}

	if mod&glfw.ModShift != 0 {
		modifiers |= ModShift
	}

	if mod&glfw.ModSuper != 0 {
		modifiers |= ModSuper
	}

	return
}

func glfwInit(opt Options) (*glfwWindow, error) {
	var wnd *glfwWindow
	var err error
	mainthread.Call(func() {
		if err = glfw.Init(); err != nil {
			return
		} else if wnd, err = newGlfwWindow(opt); err != nil {
			return
		}
		err = gl.Init(nil)
	})
	return wnd, err
}
