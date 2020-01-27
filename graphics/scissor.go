package graphics

import (
	"image"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

// Scissor represents a rectangular area of the screen
// beyond which all draw calls are clipped.
// Scissors are internally placed in a stack. Begin() and End() pairs must
// always be balanced for the stack to work properly.
//
// Scissor is a type alias of image.Rectangle and can be created as follows:
//  graphics.Scissor(image.Rect(x0, y0, x1, y1))
//  graphics.Scissor(image.Point{x0, y0}, image.Point{x1, y1})
//
// To apply a Scissor:
//  scissor.Begin()
//  ... (anything outside scissor is clipped)
//  scissor.End()
//
// The special scissor ZeroScissor disables scissoring:
//  graphics.ZeroScissor.Begin()
//  ... (no clipping here)
//  graphics.ZeroScissor.End()
type Scissor image.Rectangle

type scissorStack struct {
	stack   []Scissor
	current Scissor
}

var scissorstack = &scissorStack{}

// ZeroScissor represents the absence of clipping.
var ZeroScissor = Scissor(image.Rectangle{})

// Applies the scissor.
func (s Scissor) Begin() {
	scissorstack.push(s)
}

// Re-applies the previous scissor.
func (s Scissor) End() {
	scissorstack.pop()
}

func (ss *scissorStack) setScissor(scissor Scissor) {
	if scissor != ss.current {
		if scissor == ZeroScissor {
			gl.Disable(gl.SCISSOR_TEST)
		} else if ss.current == ZeroScissor {
			gl.Enable(gl.SCISSOR_TEST)
			gl.Scissor(image.Rectangle(scissor))
		} else {
			gl.Scissor(image.Rectangle(scissor))
		}
		ss.current = scissor
	}
}

func (ss *scissorStack) push(next Scissor) {
	ss.stack = append(ss.stack, next)
	ss.setScissor(next)
}

func (ss *scissorStack) pop() {
	prev := ss.stack[len(ss.stack)-1]
	ss.stack = ss.stack[:len(ss.stack)-1]
	ss.setScissor(prev)
}
