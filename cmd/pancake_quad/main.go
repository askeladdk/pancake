package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/graphics/postprocessing"
	"github.com/go-gl/mathgl/mgl32"
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

var vertices = []float32{
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

const Tau = 2 * math.Pi

var LogicalResolution = image.Point{640, 480}

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

type quadState struct {
	program *graphics.ShaderProgram
	vslice  *graphics.IndexedVertexSlice
	texture *graphics.Texture
	frame   *graphics.Framebuffer
	pp      *postprocessing.Processor
}

func (state *quadState) Begin(loop pancake.Loop) {
	fmt.Println(gl.GetString(gl.VERSION))

	if program, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		panic(err)
	} else if img, err := loadImage("gamer-gopher.png"); err != nil {
		panic(err)
	} else if frame, err := graphics.NewFramebuffer(LogicalResolution, graphics.FilterNearest, false); err != nil {
		panic(err)
	} else {
		// drawing variables
		state.program = program
		state.texture = graphics.NewTexture(
			img.Bounds().Size(), graphics.FilterLinear, graphics.ColorFormatRGBA, img.Pix)
		buffer := graphics.NewBuffer(quadFormat, 6, vertices)
		ebo := graphics.NewIndexBufferUint8(indices)
		state.vslice = graphics.NewIndexedVertexSlice(ebo, buffer)

		// post processing variables
		state.pp = postprocessing.NewProcessor()
		state.frame = frame

		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	}
}

func (state *quadState) End(loop pancake.Loop) {
}

func (state *quadState) Frame(loop pancake.Loop) {
}

func (state *quadState) Draw(loop pancake.Loop) {
	gl.Viewport(state.frame.Bounds())
	state.frame.Begin()

	gl.ClearColor(1, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// setup modelview matrix
	// scale to the size of the texture
	// and translate to centre at the frame
	texsz := state.texture.Size()
	framesz := state.frame.Bounds().Size()
	scale := mgl32.Scale2D(float32(texsz.X), float32(texsz.Y))
	translate := mgl32.Translate2D(float32(framesz.X/2), float32(framesz.Y/2))
	modelview := translate.Mul3(scale)

	projection := mgl32.Ortho2D(
		0,
		float32(LogicalResolution.X),
		float32(LogicalResolution.Y),
		0,
	)

	state.texture.Begin()
	state.program.Begin()
	state.program.SetUniform("projection", projection)
	state.program.SetUniform("modelview", modelview)
	state.vslice.Begin()
	state.vslice.Draw(gl.TRIANGLES)
	state.vslice.End()
	state.program.End()
	state.texture.End()

	state.frame.End()

	state.pp.Do(
		postprocessing.Targets{
			state.frame,
			loop.Window().Framebuffer(),
		},
		postprocessing.Effects{
			state.pp.Blit(),
		},
	)
}

func main() {
	loop := pancake.NewFixedTimeStepLoop(&quadState{}, 60)
	pancake.Run(pancake.Options{
		WindowSize: image.Point{800, 600},
		Resolution: image.Point{640, 480},
		Title:      "hello gopher",
	}, loop.Run)
}
