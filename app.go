package pancake

import (
	"image"
	"time"

	gl "github.com/askeladdk/pancake/opengl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.design/x/mainthread"
)

type constError string

func (s constError) Error() string {
	return string(s)
}

// ErrQuit signals that the event loop must end.
var ErrQuit error = constError("quit application")

func makeInputFlags(action glfw.Action, mod glfw.ModifierKey) Modifiers {
	var flags Modifiers
	if action == glfw.Press {
		flags |= ModPressed
	} else if action == glfw.Release {
		flags |= ModReleased
	} else if action == glfw.Repeat {
		flags |= ModRepeated
	}

	if mod&glfw.ModAlt != 0 {
		flags |= ModAlt
	}

	if mod&glfw.ModControl != 0 {
		flags |= ModControl
	}

	if mod&glfw.ModShift != 0 {
		flags |= ModShift
	}

	if mod&glfw.ModSuper != 0 {
		flags |= ModSuper
	}

	return flags
}

func makeWindow(opt Options) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	wnd, err := glfw.CreateWindow(opt.WindowSize.X, opt.WindowSize.Y, opt.Title, nil, nil)
	if err != nil {
		return nil, err
	}

	wnd.MakeContextCurrent()
	return wnd, nil
}

type App interface {
	Scissor(r image.Rectangle) Scissor
	Resolution() image.Point
	Framebuffer() *Framebuffer
	FrameRate() int
	SetTitle(string)
	Events(func(interface{}) error) error
	Begin()
	End()
}

type app struct {
	windowScale   int
	resolution    image.Point
	viewport      image.Rectangle
	window        *glfw.Window
	framebuffer   *Framebuffer
	inputEvents   []interface{}
	mousePosition image.Point
	deltaTime     float64
	frameRate     int
	cursorEntered bool
}

func (app *app) Scissor(r image.Rectangle) Scissor {
	vpsz := app.viewport.Size()
	scale := vpsz.X / app.resolution.X
	r.Min = r.Min.Mul(scale)
	r.Max = r.Max.Mul(scale)
	r.Min.Y, r.Max.Y = vpsz.Y-r.Max.Y, vpsz.Y-r.Min.Y
	return Scissor(r)
}

func (app *app) Framebuffer() *Framebuffer {
	return app.framebuffer
}

func (app *app) FrameRate() int {
	return app.frameRate
}

func (app *app) SetTitle(title string) {
	mainthread.Call(func() {
		app.window.SetTitle(title)
	})
}

func (app *app) Events(eventh func(interface{}) error) error {
	// loop regulator variables
	deltaTime := app.deltaTime
	accumulator := float64(0)
	t0 := time.Now()

	// frame counter variables
	frameRate := 0
	ft0 := time.Now()

mainloop:
	for {
		t1 := time.Now()
		accumulator += t1.Sub(t0).Seconds()
		t0 = t1

		for accumulator >= deltaTime {
			accumulator -= deltaTime

			if app.window.ShouldClose() {
				app.window.SetShouldClose(false)
				app.inputEvents = append(app.inputEvents, QuitEvent{})
			}

			mainthread.Call(glfw.PollEvents)

			for _, inputEvent := range app.inputEvents {
				if err := eventh(inputEvent); err == ErrQuit {
					break mainloop
				} else if err != nil {
					return err
				}
			}
			app.inputEvents = app.inputEvents[:0]

			if err := eventh(FrameEvent{deltaTime}); err == ErrQuit {
				break mainloop
			} else if err != nil {
				return err
			}

			// frame counter
			frameRate++
			ft1 := time.Now()
			if ft1.Sub(ft0).Seconds() >= 1 {
				app.frameRate, frameRate, ft0 = frameRate, 0, ft1
			}
		}

		alpha := accumulator / deltaTime

		if err := eventh(DrawEvent{alpha}); err == ErrQuit {
			break mainloop
		} else if err != nil {
			return err
		}

		mainthread.Call(app.window.SwapBuffers)
	}

	return nil
}

func (app *app) Begin() {
	app.framebuffer.Begin()
	gl.Viewport(app.framebuffer.Bounds())
}

func (app *app) End() {
	app.framebuffer.End()
	screen := Framebuffer{}
	screen.Begin()
	gl.Viewport(app.viewport)
	screen.End()
	app.framebuffer.BlitTo(&screen,
		app.framebuffer.Bounds(), app.viewport,
		gl.COLOR_BUFFER_BIT, FilterLinear)
}

func (app *app) Resolution() image.Point {
	return app.resolution
}

func (app *app) charCallback(_ *glfw.Window, char rune) {
	app.inputEvents = append(app.inputEvents, CharEvent{
		Char: char,
	})
}

func (app *app) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mod glfw.ModifierKey) {
	app.inputEvents = append(app.inputEvents, KeyEvent{
		Key:       Key(key),
		Modifiers: makeInputFlags(action, mod),
		Scancode:  scancode,
	})
}

func (app *app) cursorEnterCallback(_ *glfw.Window, entered bool) {
	app.cursorEntered = entered
}

func (app *app) cursorCallback(_ *glfw.Window, x, y float64) {
	mouse := image.Point{int(x), int(y)}.Mul(app.windowScale)

	if app.cursorEntered && mouse.In(app.viewport) {
		// Scale the mouse position from window to resolution coordinates.
		mouse = mouse.Sub(app.viewport.Min)
		vpsz := app.viewport.Size()
		app.mousePosition = image.Point{
			mouse.X * app.resolution.X / vpsz.X,
			mouse.Y * app.resolution.Y / vpsz.Y,
		}

		app.inputEvents = append(app.inputEvents, MouseMoveEvent{
			Position: app.mousePosition,
		})
	}
}

func (app *app) mouseCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if app.cursorEntered {
		app.inputEvents = append(app.inputEvents, MouseEvent{
			Button:    MouseButton(button),
			Modifiers: makeInputFlags(action, mod),
			Position:  app.mousePosition,
		})
	}
}

type Options struct {
	WindowSize image.Point
	Resolution image.Point
	Title      string
	FrameRate  int
}

func Main(opt Options, run func(App) error) error {
	if opt.Resolution == (image.Point{}) {
		opt.Resolution = opt.WindowSize
	}

	if opt.FrameRate <= 0 {
		opt.FrameRate = 60
	}

	var err error
	mainthread.Init(func() {
		var window *glfw.Window

		if window, err = initglfw(opt); err != nil {
			return
		}

		defer func() {
			mainthread.Call(glfw.Terminate)
		}()

		w, h := window.GetFramebufferSize()
		viewport := logicalViewport(image.Point{w, h}, opt.Resolution)
		framebuffer, _ := NewFramebuffer(viewport.Size(), FilterLinear, true)

		a := app{
			windowScale: w / opt.WindowSize.X,
			window:      window,
			deltaTime:   1 / float64(opt.FrameRate),
			viewport:    viewport,
			resolution:  opt.Resolution,
			framebuffer: framebuffer,
		}

		window.SetKeyCallback(a.keyCallback)
		window.SetCharCallback(a.charCallback)
		window.SetCursorPosCallback(a.cursorCallback)
		window.SetCursorEnterCallback(a.cursorEnterCallback)
		window.SetMouseButtonCallback(a.mouseCallback)

		gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer(0))

		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

		err = run(&a)
	})

	return err
}

func initglfw(opt Options) (*glfw.Window, error) {
	var window *glfw.Window
	var err error
	mainthread.Call(func() {
		if err = glfw.Init(); err != nil {
			return
		} else if window, err = makeWindow(opt); err != nil {
			return
		}
		err = gl.Init(nil)
	})
	return window, err
}
