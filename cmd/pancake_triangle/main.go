package main

import (
	"fmt"
	"image"
	"time"

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

var triangleFormat = graphics.AttribFormat{
	graphics.Vec2, // x, y
	graphics.Vec3, // r, g, b
}

func run(win pancake.Window) {
	fmt.Println(gl.GetString(gl.VERSION))

	if sh, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		panic(err)
	} else {
		buffer := graphics.NewBuffer(triangleFormat, 6, triangle)
		vx := graphics.NewVertexSlice(buffer)

		t0 := time.Now()
		for !win.ShouldClose() {
			gl.ClearColor(1, 0, 0, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT)

			dt := time.Now().Sub(t0).Seconds()

			v := vx.Slice(int(dt)%4, int(dt)%4+3)

			sh.Begin()
			sh.SetUniform("dt", float32(dt))
			v.Begin()
			v.Draw(gl.TRIANGLES)
			v.End()
			sh.End()

			win.Update()
		}
	}
}

func main() {
	pancake.Run(pancake.WindowOptions{
		Size:  image.Point{640, 400},
		Title: "hello triangle",
	}, run)
}
