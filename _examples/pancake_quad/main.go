package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/input"
	"github.com/askeladdk/pancake/mathx"
)

var vshader = `
#version 330 core

layout(location = 0) in vec2 in_Position;
layout(location = 1) in vec2 in_Texture;

out vec2 texcoord;

uniform mat4 projection;
uniform mat3 modelview;

void main()
{
	texcoord = in_Texture;
    gl_Position = projection * vec4(modelview * vec3(in_Position, 1.0), 1.0);
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

var vertices = []float64{
	// x, y, u, v
	// -1, -1, 0, 1,
	// +1, -1, 1, 1,
	// +1, +1, 1, 0,
	// -1, +1, 0, 0,

	-.5, -.5, 0, 0,
	-.5, +.5, 0, 1,
	+.5, -.5, 1, 0,
	+.5, +.5, 1, 1,
}

// var indices = []uint8{0, 1, 2, 0, 2, 3}

var indices = []uint8{0, 1, 2, 1, 2, 3}

var quadFormat = graphics.AttribFormat{
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

func run(app pancake.App) error {
	fmt.Println(gl.GetString(gl.VERSION))

	if program, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		return err
	} else if img, err := loadImage("gamer-gopher.png"); err != nil {
		return err
	} else {
		// drawing variables
		texture := graphics.NewTexture(
			img.Bounds().Size(), graphics.FilterLinear, graphics.ColorFormatRGBA, img.Pix)
		buffer := graphics.NewBuffer(quadFormat, 6, vertices)
		ebo := graphics.NewIndexBufferUint8(indices)
		vslice := graphics.NewIndexedVertexSlice(ebo, buffer)

		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		return app.Events(func(ev interface{}) error {
			switch e := ev.(type) {
			case pancake.QuitEvent:
				return pancake.Quit
			case pancake.KeyEvent:
				if e.Flags.Pressed() && e.Key == input.KeyEscape {
					return pancake.Quit
				}
			case pancake.DrawEvent:
				app.Begin()
				gl.ClearColor(1, 0, 0, 1)
				gl.Clear(gl.COLOR_BUFFER_BIT)

				// setup modelview matrix
				// scale to the size of the texture
				// and translate to centre at the frame
				texsz := texture.Size()
				framesz := app.Resolution()
				modelview := mathx.
					ScaleAff3(mathx.FromPoint(texsz)).
					Translated(mathx.FromPoint(framesz).Mul(0.5))

				projection := mathx.Ortho2D(
					0,
					float64(framesz.X),
					float64(framesz.Y),
					0,
				)

				texture.Begin()
				program.Begin()
				program.SetUniform("projection", projection)
				program.SetUniform("modelview", modelview)
				vslice.Begin()
				vslice.Draw(gl.TRIANGLES)
				vslice.End()
				program.End()
				texture.End()

				app.End()
			}

			return nil
		})
	}
}

func main() {
	opt := pancake.Options{
		WindowSize: image.Point{800, 600},
		Resolution: image.Point{640, 480},
		Title:      "hello gopher",
		FrameRate:  60,
	}

	pancake.Main(opt, run)
}
