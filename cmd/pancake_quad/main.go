package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/desktop"
	"github.com/askeladdk/pancake/graphics"
	"github.com/faiface/mainthread"
)

var vshader = `
#version 330 core

layout(location = 0) in vec2 in_Position;
layout(location = 1) in vec2 in_Texture;

out vec2 texcoord;

void main()
{
	texcoord = in_Texture;
    gl_Position = vec4(in_Position, 0.0, 1.0);
}
`

var fshader = `
#version 330 core

in vec2 texcoord;

out vec4 FragColor;

uniform sampler2D tex;

void main()
{
    FragColor = texture(tex, texcoord);
}
`

var quad = []float32{
	-1, -1, 0, 1,
	+1, -1, 1, 1,
	+1, +1, 1, 0,
	-1, -1, 0, 1,
	+1, +1, 1, 0,
	-1, +1, 0, 0,
}

var quadFormat = graphics.AttrFormat{
	graphics.Vec2, // x, y
	graphics.Vec2, // u, v
}

func loadImage(path string) (*image.NRGBA, error) {
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else if img, err := png.Decode(f); err != nil {
		return nil, err
	} else {
		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, image.Point{0, 0}, draw.Src)
		return nrgba, nil
	}
}

func run(win pancake.Window) {
	var err error
	var vx graphics.VertexSlice
	var sh graphics.Shader
	var driver graphics.Driver
	var tx graphics.Texture

	mainthread.Call(func() {
		driver = graphics.Get()

		fmt.Println(driver.Version())

		if vx, err = driver.NewVertexSlice(quadFormat, 6, quad); err != nil {
			panic(err)
		} else if sh, err = driver.NewShader(vshader, fshader); err != nil {
			panic(err)
		} else if img, err := loadImage("gamer-gopher.png"); err != nil {
			panic(err)
		} else if tx, err = driver.NewTexture(
			img.Bounds().Size(), graphics.FilterLinear,
			graphics.ColorFormatRGBA, img.Pix); err != nil {
			panic(err)
		}
	})

	for !win.ShouldClose() {
		mainthread.Call(func() {
			driver.Clear(1, 0, 0, 0)

			tx.Begin()
			sh.Begin()
			vx.Begin()
			vx.Draw()
			vx.End()
			sh.End()
			tx.End()
		})

		win.Update()
	}
}

func main() {
	desktop.Run(pancake.WindowOptions{
		Width:  640,
		Height: 400,
		Title:  "hello quad",
	}, run)
}
