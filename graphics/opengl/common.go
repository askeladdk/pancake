package opengl

import (
	"errors"
)

type (
	Attrib      uint32
	Buffer      uint32
	Enum        uint32
	Framebuffer uint32
	Program     uint32
	Shader      uint32
	Texture     uint32
	Uniform     int32
	VertexArray uint32
)

var (
	ErrNotSupported = errors.New("function not supported")
)
