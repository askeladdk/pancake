package main

import (
	"fmt"
	"time"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/desktop"
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

uniform float dt;

void main()
{
    out_color = vec4(color * sin(dt), 1.0);
}
`

var triangle = []float32{
	// x, y, r, g, b
	+0.0, +0.5, 1, 0, 0,
	+0.5, -0.5, 0, 1, 0,
	-0.5, -0.5, 0, 0, 1,

	+0.0, +0.2, 1, 0, 0,
	+0.2, -0.2, 0, 1, 0,
	-0.2, -0.2, 0, 0, 1,
}

var triangleFormat = graphics.AttrFormat{
	graphics.Vec2, // x, y
	graphics.Vec3, // r, g, b
}

func run(win pancake.Window) {
	fmt.Println(gl.GetString(gl.VERSION))

	if vx, err := graphics.NewVertexSlice(triangleFormat, 6, triangle); err != nil {
		panic(err)
	} else if sh, err := graphics.NewShader(vshader, fshader); err != nil {
		panic(err)
	} else {
		t0 := time.Now()
		for !win.ShouldClose() {
			gl.ClearColor(1, 0, 0, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT)

			dt := time.Now().Sub(t0).Seconds()

			v := vx.Slice(int(dt)%4, int(dt)%4+3)

			sh.Begin()
			sh.SetUniform("dt", float32(dt))
			v.Begin()
			v.Draw()
			v.End()
			sh.End()

			win.Update()
		}
	}
}

func main() {
	desktop.Run(pancake.WindowOptions{
		Width:  640,
		Height: 400,
		Title:  "hello triangle",
	}, run)
}
