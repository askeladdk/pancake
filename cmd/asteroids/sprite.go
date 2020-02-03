package main

import (
	"image"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/go-gl/mathgl/mgl32"
)

const vshader = `
#version 330 core

layout(location = 0) in vec2 in_Position;
layout(location = 1) in vec2 in_Texture;

out vec2 texcoord;

uniform mat4 projection;
uniform mat3 modelview;
uniform mat3 uvmat;

void main()
{
	texcoord = (uvmat * vec3(in_Texture, 1.0)).xy;
    gl_Position = projection * vec4(modelview * vec3(in_Position, 1.0), 1.0);
}
`

const fshader = `
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
	// x, y, u, v
	-.5, -.5, 0, 0,
	-.5, +.5, 0, 1,
	+.5, -.5, 1, 0,

	-.5, +.5, 0, 1,
	+.5, -.5, 1, 0,

	+.5, +.5, 1, 1,
}

// var indices = []uint8{0, 1, 2, 0, 2, 3}

var quadFormat = graphics.AttribFormat{
	graphics.Vec2, // x, y
	graphics.Vec2, // u, v
}

type SpriteDrawer struct {
	projection mgl32.Mat4
	program    *graphics.ShaderProgram
	vbuffer    *graphics.Buffer
	vslice     *graphics.VertexSlice
}

func NewSpriteDrawer(resolution image.Point) (*SpriteDrawer, error) {
	if program, err := graphics.NewShaderProgram(vshader, fshader); err != nil {
		return nil, err
	} else {
		vbuffer := graphics.NewBuffer(quadFormat, 6, quad)
		vslice := graphics.NewVertexSlice(vbuffer)
		return &SpriteDrawer{
			projection: mgl32.Ortho2D(
				0,
				float32(resolution.X),
				float32(resolution.Y),
				0,
			),
			program: program,
			vbuffer: vbuffer,
			vslice:  vslice,
		}, nil
	}
}

func (sd *SpriteDrawer) Begin() {
	sd.program.Begin()
	sd.program.SetUniform("projection", sd.projection)
}

func (sd *SpriteDrawer) DrawImage(img *Image, position mgl32.Vec2) {
	modelview := mgl32.Mat3{
		img.Size[0], 0, 0,
		0, img.Size[1], 0,
		position[0], position[1], 1,
	}
	img.Texture.Begin()
	sd.program.SetUniform("modelview", modelview)
	sd.program.SetUniform("uvmat", img.UVBounds)
	sd.vslice.Begin()
	sd.vslice.Draw(gl.TRIANGLES)
	sd.vslice.End()
	img.Texture.End()
}

func (sd *SpriteDrawer) End() {
	sd.program.End()
}
