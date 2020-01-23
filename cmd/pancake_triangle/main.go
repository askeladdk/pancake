package main

import (
	"fmt"
	"image"

	"github.com/askeladdk/pancake/input"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var vshader = `
#version 330 core

layout(location = 0) in vec2 in_position;
layout(location = 1) in vec3 in_color;

out vec3 color;

void main()
{
	color = in_color;
    gl_Position = vec4(in_position, 0.0, 1.0);
}
`

var fshader = `
#version 330 core

in vec3 color;
out vec4 out_color;

uniform float t;

void main()
{
    out_color = vec4(color * sin(t), 1.0);
}
`

var triangle = []float32{
	// x, y, r, g, b
	+0.0, +0.5, 1, 0, 0,
	+0.5, -0.5, 0, 1, 0,
	-0.5, -0.5, 0, 0, 1,
}

var triangleFormat = graphics.AttribFormat{
	graphics.Vec2, // x, y
	graphics.Vec3, // r, g, b
}

type triangleState struct {
	program     *graphics.ShaderProgram
	buffer      *graphics.Buffer
	vslice      *graphics.VertexSlice
	tn, tl      float64
	interpolate bool
	escape      bool
}

func (state *triangleState) KeyEvent(event input.KeyEvent) {
	if event.Flags.Pressed() {
		switch event.Key {
		case input.KeySpace:
			state.interpolate = !state.interpolate
		case input.KeyEscape:
			state.escape = true
		}
	}
}

func (state *triangleState) Begin(loop pancake.Loop) {
	loop.Window().SetKeyEventHandler(state)

	fmt.Println(gl.GetString(gl.VERSION))
	if program, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		panic(err)
	} else {
		state.program = program
	}

	state.buffer = graphics.NewBuffer(triangleFormat, 3, triangle)
	state.vslice = graphics.NewVertexSlice(state.buffer)
}

func (state *triangleState) End(loop pancake.Loop) {
	fmt.Println("So long folks")
}

func (state *triangleState) Frame(loop pancake.Loop) {
	if state.escape {
		loop.Transition(pancake.EmptyState)
		return
	}

	// store the last state and calculate the next state.
	state.tl, state.tn = state.tn, state.tn+loop.DeltaTime()

	loop.Window().SetTitle(fmt.Sprintf("FPS: %d | Elapsed: %.1fs | SPACE toggles interpolation", loop.FrameRate(), state.tn))
}

func (state *triangleState) Draw(loop pancake.Loop) {
	gl.ClearColor(1, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	var t float64
	if state.interpolate {
		// linear interpolation between current and previous state.
		a := loop.Alpha()
		t = state.tn*a + state.tl*(1-a)
	} else {
		// no interpolation
		t = state.tn
	}

	state.program.Begin()
	state.program.SetUniform("t", float32(t))
	state.vslice.Begin()
	state.vslice.Draw(gl.TRIANGLES)
	state.vslice.End()
	state.program.End()
}

func main() {
	// targetFrameRate is deliberately set low in order
	// to demonstrate interpolation in Draw().
	loop := pancake.NewFixedTimeStepLoop(&triangleState{}, 5)
	pancake.Run(pancake.Options{
		WindowSize: image.Point{640, 400},
		Resolution: image.Point{640, 400},
		Title:      "hello triangle",
	}, loop.Run)
}
