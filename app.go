package pancake

import (
	"context"
	"image"
	"time"

	gl "github.com/askeladdk/pancake/opengl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.design/x/mainthread"
)

// App is the window and event handling context.
type App interface {
	// Begin should be called when a draw event begins.
	Begin()

	// End should be called when a draw even ends.
	End()

	// Events returns the events channel.
	Events() <-chan interface{}

	// FrameRate reports the current framerate.
	FrameRate() int

	// Scissor creates a Scissor that is scaled to the viewport.
	Scissor(image.Rectangle) Scissor

	// Bounds reports the bounds of the viewport.
	Bounds() image.Rectangle

	// SetTitle sets the window title.
	SetTitle(string)
}

type app struct {
	*glfwWindow
	deltaTime float64
	frameRate int
	eventch   chan interface{}
}

func (app *app) FrameRate() int { return app.frameRate }

func (app *app) Events() <-chan interface{} { return app.eventch }

func (app *app) loop(ctx context.Context) {
	// loop regulator variables
	deltaTime := app.deltaTime
	accumulator := float64(0)
	t0 := time.Now()

	// frame counter variables
	frameRate := 0
	ft0 := time.Now()

	for {
		t1 := time.Now()
		accumulator += t1.Sub(t0).Seconds()
		t0 = t1

		for accumulator >= deltaTime {
			accumulator -= deltaTime

			if !app.InputEvents(ctx, app.eventch) {
				return
			}

			select {
			case <-ctx.Done():
				return
			case app.eventch <- FrameEvent{deltaTime}:
			}

			// frame counter
			frameRate++
			ft1 := time.Now()
			if ft1.Sub(ft0).Seconds() >= 1 {
				app.frameRate, frameRate, ft0 = frameRate, 0, ft1
			}
		}

		select {
		case <-ctx.Done():
			return
		case app.eventch <- DrawEvent{accumulator / deltaTime}:
		}
	}
}

// Options specifies window options.
type Options struct {
	// WindowSize is the size of the window.
	WindowSize image.Point
	// Resolution is the logical resolution in pixels.
	// Can be smaller than WindowSize.
	Resolution image.Point
	// Title is the window title.
	Title string
	// FrameRate is the target frame rate.
	FrameRate int
}

// Main initializes the window and starts the event loop.
func Main(opt Options, run func(App) error) error {
	if opt.Resolution == (image.Point{}) {
		opt.Resolution = opt.WindowSize
	}

	if opt.FrameRate <= 0 {
		opt.FrameRate = 60
	}

	var err error
	mainthread.Init(func() {
		var window *glfwWindow

		if window, err = glfwInit(opt); err != nil {
			return
		}

		defer func() {
			mainthread.Call(glfw.Terminate)
		}()

		a := app{
			glfwWindow: window,
			deltaTime:  1 / float64(opt.FrameRate),
			eventch:    make(chan interface{}),
		}

		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			err = run(&a)
			cancel()
		}()
		a.loop(ctx)
	})

	return err
}
