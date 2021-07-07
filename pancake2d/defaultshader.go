package pancake2d

import (
	"github.com/askeladdk/pancake"
)

// DefaultVertexShader is the default vertex shader.
const DefaultVertexShader = `
#version 330 core

layout(location = 0) in vec2 in_Position;
layout(location = 1) in vec2 in_Texture;
layout(location = 2) in vec4 in_Color;
layout(location = 3) in float in_ZOrder;

out vec2 f_Texture;
out vec4 f_Color;

uniform mat4 u_Projection;

void main()
{
	f_Texture = in_Texture;
	f_Color = in_Color;
    gl_Position = u_Projection * vec4(in_Position, in_ZOrder, 1);
}
`

// DefaultFragmentShader is the default fragment shader.
const DefaultFragmentShader = `
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

var defaultShader *pancake.ShaderProgram

// DefaultShader returns the default shader program.
func DefaultShader() *pancake.ShaderProgram {
	if defaultShader != nil {
		return defaultShader
	} else if p, err := pancake.NewShaderProgram(DefaultVertexShader, DefaultFragmentShader); err != nil {
		panic(err)
	} else {
		defaultShader = p
		return p
	}
}
