package postprocessing

import (
	"image"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/go-gl/mathgl/mgl32"
)

// Pass-through vertex shader.
const VertexShader = `
#version 330 core

layout(location = 0) in vec2 position;
layout(location = 1) in vec2 texcoord;

out vec2 UV;

void main()
{
	UV = texcoord;
    gl_Position = vec4(position, 0, 1);
}
`

// Pass-through fragment shader.
const fragmentShader = `
#version 330 core

in vec2 UV;

out vec4 COLOR;

uniform sampler2D TEXTURE;

void main()
{
    COLOR = texture(TEXTURE, UV);
}
`

var (
	vertices = []float32{
		// -1, -1, 0, 1,
		// +1, -1, 1, 1,
		// +1, +1, 1, 0,
		// -1, -1, 0, 1,
		// +1, +1, 1, 0,
		// -1, +1, 0, 0,
		-1, -1, 0, 0,
		-1, +1, 0, 1,
		+1, -1, 1, 0,
		-1, +1, 0, 1,
		+1, -1, 1, 0,
		+1, +1, 1, 1,
	}

	attribFormat = graphics.AttribFormat{
		graphics.Vec2,
		graphics.Vec2,
	}
)

// Slice of target framebuffers.
type Targets []*graphics.Framebuffer

type Effects []*graphics.ShaderProgram

func rectToVec4(r image.Rectangle) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(r.Min.X),
		float32(r.Min.Y),
		float32(r.Max.X),
		float32(r.Max.Y),
	}
}

type Processor struct {
	blit   *graphics.ShaderProgram
	vslice *graphics.VertexSlice
}

func NewProcessor() *Processor {
	if blit, err := NewEffect(fragmentShader); err != nil {
		panic(err)
	} else {
		vbo := graphics.NewBuffer(attribFormat, 6, vertices)
		return &Processor{
			blit:   blit,
			vslice: graphics.NewVertexSlice(vbo),
		}
	}
}

func (proc *Processor) Blit() *graphics.ShaderProgram {
	return proc.blit
}

func (proc *Processor) applyEffect(effect *graphics.ShaderProgram, src, dst *graphics.Framebuffer) {
	gl.Viewport(dst.Bounds())
	dst.Begin()
	src.Color().Begin()
	effect.Begin()
	effect.SetUniform("TEXTURE", 0)
	effect.SetUniform("SRCBOUNDS", rectToVec4(src.Bounds()))
	effect.SetUniform("DSTBOUNDS", rectToVec4(dst.Bounds()))
	proc.vslice.Draw(gl.TRIANGLES)
	effect.End()
	src.Color().End()
	dst.End()
}

func (proc *Processor) Do(targets Targets, effects Effects) Targets {
	proc.vslice.Begin()
	defer proc.vslice.End()

	for _, effect := range effects {
		if len(targets) < 2 {
			return targets
		} else {
			proc.applyEffect(effect, targets[0], targets[1])
			targets = targets[1:]
		}
	}

	return targets
}

func NewEffect(fshader string) (*graphics.ShaderProgram, error) {
	if prg, err := graphics.NewShaderProgram(VertexShader, fshader); err != nil {
		return nil, err
	} else {
		return prg, nil
	}
}
