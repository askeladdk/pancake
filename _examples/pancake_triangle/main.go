package main

import (
	"fmt"
	"image"

	"github.com/askeladdk/pancake/input"
	"github.com/askeladdk/pancake/mathx"

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

var triangle = []float64{
	// x, y, r, g, b
	+0.0, +0.5, 1, 0, 0,
	+0.5, -0.5, 0, 1, 0,
	-0.5, -0.5, 0, 0, 1,
}

var triangleFormat = graphics.AttribFormat{
	graphics.Vec2, // x, y
	graphics.Vec3, // r, g, b
}

func run(app pancake.App) error {
	var program *graphics.ShaderProgram
	var tn, tl float64
	var interpolate bool

	if p, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		return err
	} else {
		program = p
	}

	buffer := graphics.NewBuffer(triangleFormat, 3, triangle)
	vslice := graphics.NewVertexSlice(buffer)

	return app.Events(func(ev interface{}) error {
		switch e := ev.(type) {
		case pancake.QuitEvent:
			return pancake.Quit
		case pancake.KeyEvent:
			if e.Flags.Pressed() {
				if e.Key == input.KeyEscape {
					return pancake.Quit
				} else if e.Key == input.KeySpace {
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

		return nil
	})
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
