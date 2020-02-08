package graphics2d

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/mathx"
)

const vertexShader = `
#version 330 core

layout(location = 0) in vec2 in_Position;
layout(location = 1) in vec2 in_Texture;
layout(location = 2) in vec4 in_Color;

out vec2 f_Texture;
out vec4 f_Color;

uniform mat4 u_Projection;

void main()
{
	f_Texture = in_Texture;
	f_Color = in_Color;
    gl_Position = u_Projection * vec4(in_Position, 0, 1);
}
`

const fragmentShader = `
#version 330 core

in vec2 f_Texture;
in vec4 f_Color;

out vec4 out_FragColor;

uniform sampler2D u_Texture;

void main()
{
    out_FragColor = texture(u_Texture, f_Texture) * f_Color;
}
`

var Quad = Mesh{
	Vertices: []Vertex{
		Vertex{
			XY:   mathx.Vec2{-.5, -.5},
			UV:   mathx.Vec2{0, 0},
			RGBA: color.NRGBA{255, 255, 255, 255},
		},
		Vertex{
			XY:   mathx.Vec2{-.5, +.5},
			UV:   mathx.Vec2{0, 1},
			RGBA: color.NRGBA{255, 255, 255, 255},
		},
		Vertex{
			XY:   mathx.Vec2{+.5, -.5},
			UV:   mathx.Vec2{1, 0},
			RGBA: color.NRGBA{255, 255, 255, 255},
		},
		Vertex{
			XY:   mathx.Vec2{+.5, +.5},
			UV:   mathx.Vec2{1, 1},
			RGBA: color.NRGBA{255, 255, 255, 255},
		},
	},
	Indices:  []uint32{0, 1, 2, 1, 2, 3},
	DrawMode: gl.TRIANGLES,
}

var defaultShader *graphics.ShaderProgram

func DefaultShader() *graphics.ShaderProgram {
	if defaultShader != nil {
		return defaultShader
	} else if p, err := graphics.NewShaderProgram(vertexShader, fragmentShader); err != nil {
		panic(err)
	} else {
		defaultShader = p
		return p
	}
}
