package graphics

import (
	"errors"
	"fmt"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"

	"github.com/go-gl/mathgl/mgl32"
)

var shaderBinder = newBinder(func(prog uint32) {
	gl.BindProgram(gl.Program(prog))
})

type shader struct {
	ref   gl.Program
	attrs map[string]gl.Uniform
}

func (this *shader) Begin() {
	shaderBinder.bind(uint32(this.ref))
}

func (this *shader) End() {
	shaderBinder.unbind()
}

func (this *shader) delete() {
	gl.DeleteProgram(this.ref)
}

func (this *shader) getUniformLocation(name string) gl.Uniform {
	if loc, ok := this.attrs[name]; !ok {
		loc = gl.GetUniformLocation(this.ref, name)
		this.attrs[name] = loc
		return loc
	} else {
		return loc
	}
}

func (this *shader) SetUniform(name string, value interface{}) bool {
	if loc := this.getUniformLocation(name); loc < 0 {
		return false
	} else {
		switch v := value.(type) {
		case int:
			gl.Uniform1i(loc, v)
		case uint:
			gl.Uniform1ui(loc, v)
		case float32:
			gl.Uniform1f(loc, v)
		case mgl32.Vec2:
			gl.Uniform2fv(loc, v[:])
		case mgl32.Vec3:
			gl.Uniform3fv(loc, v[:])
		case mgl32.Vec4:
			gl.Uniform4fv(loc, v[:])
		case mgl32.Mat2:
			gl.UniformMatrix2fv(loc, v[:])
		case mgl32.Mat3:
			gl.UniformMatrix3fv(loc, v[:])
		case mgl32.Mat4:
			gl.UniformMatrix4fv(loc, v[:])
		default:
			panic(errors.New("invalid type"))
		}
		return true
	}
}

func compileShaderSource(source string, xtype gl.Enum) (gl.Shader, error) {
	ref := gl.CreateShader(xtype)
	gl.ShaderSource(ref, source)
	gl.CompileShader(ref)

	if gl.GetShaderi(ref, gl.COMPILE_STATUS) == gl.FALSE {
		return 0, fmt.Errorf(gl.GetShaderInfoLog(ref))
	}

	return ref, nil
}

func NewShader(vshader, fshader string) (*shader, error) {
	if vref, err := compileShaderSource(vshader, gl.VERTEX_SHADER); err != nil {
		return nil, err
	} else if fref, err := compileShaderSource(fshader, gl.FRAGMENT_SHADER); err != nil {
		return nil, err
	} else {
		defer gl.DeleteShader(vref)
		defer gl.DeleteShader(fref)

		prog := &shader{
			ref:   gl.CreateProgram(),
			attrs: map[string]gl.Uniform{},
		}
		runtime.SetFinalizer(prog, (*shader).delete)

		gl.AttachShader(prog.ref, vref)
		gl.AttachShader(prog.ref, fref)
		gl.LinkProgram(prog.ref)

		if gl.GetProgrami(prog.ref, gl.LINK_STATUS) == gl.FALSE {
			return nil, fmt.Errorf(gl.GetProgramInfoLog(prog.ref))
		} else {
			return prog, nil
		}
	}
}
