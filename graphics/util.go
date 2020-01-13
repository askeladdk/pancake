package graphics

import (
	"fmt"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

// GL_NO_ERROR
//     No error has been recorded. The value of this symbolic constant is guaranteed to be 0.
// GL_INVALID_ENUM
//     An unacceptable value is specified for an enumerated argument. The offending command is ignored and has no other side effect than to set the error flag.
// GL_INVALID_VALUE
//     A numeric argument is out of range. The offending command is ignored and has no other side effect than to set the error flag.
// GL_INVALID_OPERATION
//     The specified operation is not allowed in the current state. The offending command is ignored and has no other side effect than to set the error flag.
// GL_INVALID_FRAMEBUFFER_OPERATION
//     The framebuffer object is not complete. The offending command is ignored and has no other side effect than to set the error flag.
// GL_OUT_OF_MEMORY
//     There is not enough memory left to execute the command. The state of the GL is undefined, except for the state of the error flags, after this error is recorded.
// GL_STACK_UNDERFLOW
//     An attempt has been made to perform an operation that would cause an internal stack to underflow.
// GL_STACK_OVERFLOW
//     An attempt has been made to perform an operation that would cause an internal stack to overflow.
// GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT
//     Not all framebuffer attachment points are framebuffer attachment complete. This means that at least one attachment point with a renderbuffer or texture attached has its attached object no longer in existence or has an attached image with a width or height of zero, or the color attachment point has a non-color-renderable image attached, or the depth attachment point has a non-depth-renderable image attached, or the stencil attachment point has a non-stencil-renderable image attached.
//     Color-renderable formats include GL_RGBA4, GL_RGB5_A1, and GL_RGB565. GL_DEPTH_COMPONENT16 is the only depth-renderable format. GL_STENCIL_INDEX8 is the only stencil-renderable format.
// GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT
//     No images are attached to the framebuffer.
// GL_FRAMEBUFFER_UNSUPPORTED
//     The combination of internal formats of the attached images violates an implementation-dependent set of restrictions.
func errorToString(code gl.Enum) string {
	switch code {
	case gl.INVALID_ENUM:
		return "GL_INVALID_ENUM"
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return "GL_INVALID_FRAMEBUFFER_OPERATION"
	case gl.INVALID_OPERATION:
		return "GL_INVALID_OPERATION"
	case gl.INVALID_VALUE:
		return "GL_INVALID_VALUE"
	case gl.NO_ERROR:
		return "GL_NO_ERROR"
	case gl.OUT_OF_MEMORY:
		return "GL_OUT_OF_MEMORY"
	case gl.STACK_OVERFLOW:
		return "GL_STACK_OVERFLOW"
	case gl.STACK_UNDERFLOW:
		return "GL_STACK_UNDERFLOW"
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return "GL_FRAMEBUFFER_INCOMPLETE_ATTACHMENT"
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return "GL_FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT"
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return "GL_FRAMEBUFFER_UNSUPPORTED"
	default:
		panic(fmt.Errorf("unrecognised error %v", code))
	}
}

func checkError() error {
	if code := gl.GetError(); code != gl.NO_ERROR {
		return fmt.Errorf(errorToString(code))
	} else {
		return nil
	}
}

func panicError() {
	if err := checkError(); err != nil {
		panic(err)
	}
}
