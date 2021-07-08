package main

import (
	"fmt"
	"image"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/mathx"
	gl "github.com/askeladdk/pancake/opengl"
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

var triangle = []float64{
	// x, y, r, g, b
	+0.0, +0.5, 1, 0, 0,
	+0.5, -0.5, 0, 1, 0,
	-0.5, -0.5, 0, 0, 1,
}

var triangleFormat = pancake.AttribFormat{
	pancake.AttribVec2, // x, y
	pancake.AttribVec3, // r, g, b
}

func run(app pancake.App) error {
	var program *pancake.ShaderProgram
	var tn, tl float64
	var interpolate bool
	var err error

	if program, err = pancake.NewShaderProgram(vshader, fshader); err != nil {
		return err
	}

	buffer := pancake.NewVertexBuffer(triangleFormat, 3, triangle)
	vslice := pancake.NewVertexArraySlice(buffer)

	for {
		switch e := (<-app.Events()).(type) {
		case pancake.QuitEvent:
			return nil
		case pancake.KeyEvent:
			if e.Modifiers.Pressed() {
				if e.Key == pancake.KeyEscape {
					return nil
				} else if e.Key == pancake.KeySpace {
					interpolate = !interpolate
				}
			}
		case pancake.FrameEvent:
			tl, tn = tn, tn+e.DeltaTime
			app.SetTitle(fmt.Sprintf("FPS: %d | Elapsed: %.1fs | SPACE toggles interpolation", app.FrameRate(), tn))
		case pancake.DrawEvent:
			app.Begin()
			gl.ClearColor(1, 0, 0, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT)

			var t float64

			// linear interpolation between current and previous frame.
			if interpolate {
				t = mathx.Lerp(tl, tn, e.Alpha)
			} else {
				t = tn
			}

			program.Begin()
			program.SetUniform("t", t)
			vslice.Begin()
			vslice.Draw(gl.TRIANGLES)
			vslice.End()
			program.End()
			app.End()
		}
	}
}

func main() {
	opt := pancake.Options{
		WindowSize: image.Point{640, 400},
		FrameRate:  5,
	}

	if err := pancake.Main(opt, run); err != nil {
		panic(err)
	}
}
