package pancake

import (
	"time"
)

// Loop regulates frame rates and handles top-level state management.
type Loop interface {
	// Run starts the loop. It should be used as an argument to pancake.Run().
	Run(Window)

	// Transition tells the Loop to transition to a new state.
	// Pass EmptyState to terminate the Loop.
	Transition(state State)

	// Window returns the window.
	Window() Window

	// FrameRate returns the last measured frames per second.
	FrameRate() int

	// Alpha returns a value between 0 and 1 that represents
	// the remaining frame time that could not be simulated.
	// Interpolate the current and previous states with alpha
	// to smooth visual stuttering.
	//
	//  x_draw = x_current * alpha + x_previous * (1 - alpha)
	Alpha() float64

	// DeltaTime is the time in seconds between frame updates.
	DeltaTime() float64
}

// State represents a game-state that is regulated by a Loop.
type State interface {
	// Begin is called whenever the Loop transitions to the state.
	Begin(Loop)

	// Begin is called whenever the Loop transitions away from the state.
	End(Loop)

	// Frame is called whenever the state must be updated.
	Frame(Loop)

	// Draw is called whenever rendering should take place.
	Draw(Loop)
}

type emptyState struct{}

func (state *emptyState) Begin(loop Loop) {}
func (state *emptyState) End(loop Loop)   {}
func (state *emptyState) Frame(loop Loop) {}
func (state *emptyState) Draw(loop Loop)  {}

// EmptyState does nothing. Transitioning to EmptyState quits the Loop.
var EmptyState State = &emptyState{}

type fixedTimeStepLoop struct {
	state     State
	nextState State
	window    Window
	deltaTime float64
	alpha     float64
	frameRate int
}

func (loop *fixedTimeStepLoop) Transition(state State) {
	loop.nextState = state
}

func (loop *fixedTimeStepLoop) Window() Window {
	return loop.window
}

func (loop *fixedTimeStepLoop) FrameRate() int {
	return loop.frameRate
}

func (loop *fixedTimeStepLoop) Alpha() float64 {
	return loop.alpha
}

func (loop *fixedTimeStepLoop) DeltaTime() float64 {
	return loop.deltaTime
}

func (loop *fixedTimeStepLoop) Run(window Window) {
	// loop regulator variables
	t0 := time.Now()
	accumulator := float64(0)

	// frame counter variables
	frameRate := 0
	ft0 := time.Now()

	loop.window = window

mainloop:
	for {
		if window.ShouldClose() {
			break mainloop
		}

		t1 := time.Now()
		accumulator += t1.Sub(t0).Seconds()
		t0 = t1

		for accumulator >= loop.deltaTime {
			accumulator -= loop.deltaTime
			window.PollEvents()
			loop.state.Frame(loop)

			// state transition
			if loop.nextState == EmptyState {
				break mainloop
			} else if loop.nextState != nil {
				loop.state.End(loop)
				loop.nextState.Begin(loop)
				loop.state, loop.nextState = loop.nextState, nil
			}

			// frame counter
			frameRate += 1
			ft1 := time.Now()
			if ft1.Sub(ft0).Seconds() >= 1 {
				loop.frameRate, frameRate, ft0 = frameRate, 0, ft1
			}
		}

		loop.alpha = accumulator / loop.deltaTime
		loop.state.Draw(loop)
		window.SwapBuffers()
	}

	loop.state.End(loop)
}

// NewFixedTimeStepLoop creates a Loop that renders as often as possible
// while updating the logic at a fixed time interval. This is the only
// way to create a stable simulation.
// Thus DeltaTime() always returns the same value which is 1 / targetFrameRate.
func NewFixedTimeStepLoop(initialState State, targetFrameRate int) Loop {
	return &fixedTimeStepLoop{
		state:     EmptyState,
		nextState: initialState,
		deltaTime: 1 / float64(targetFrameRate),
	}
}
