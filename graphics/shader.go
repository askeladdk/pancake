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

type ShaderProgram struct {
	id    gl.Program
	attrs map[string]gl.Uniform
}

func (prg *ShaderProgram) Begin() {
	shaderBinder.bind(uint32(prg.id))
}

func (prg *ShaderProgram) End() {
	shaderBinder.unbind()
}

func (prg *ShaderProgram) delete() {
	gl.DeleteProgram(prg.id)
}

func (prg *ShaderProgram) getUniformLocation(name string) gl.Uniform {
	if loc, ok := prg.attrs[name]; !ok {
		loc = gl.GetUniformLocation(prg.id, name)
		prg.attrs[name] = loc
		return loc
	} else {
		return loc
	}
}

func (prg *ShaderProgram) SetUniform(name string, value interface{}) bool {
	if loc := prg.getUniformLocation(name); loc < 0 {
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
			gl.Uniform2fv(loc, []mgl32.Vec2{v})
		case mgl32.Vec3:
			gl.Uniform3fv(loc, []mgl32.Vec3{v})
		case mgl32.Vec4:
			gl.Uniform4fv(loc, []mgl32.Vec4{v})
		case mgl32.Mat2:
			gl.UniformMatrix2fv(loc, []mgl32.Mat2{v})
		case mgl32.Mat3:
			gl.UniformMatrix3fv(loc, []mgl32.Mat3{v})
		case mgl32.Mat4:
			gl.UniformMatrix4fv(loc, []mgl32.Mat4{v})
		case []float32:
			gl.Uniform1fv(loc, v)
		case []mgl32.Vec2:
			gl.Uniform2fv(loc, v)
		case []mgl32.Vec3:
			gl.Uniform3fv(loc, v)
		case []mgl32.Vec4:
			gl.Uniform4fv(loc, v)
		case []mgl32.Mat2:
			gl.UniformMatrix2fv(loc, v)
		case []mgl32.Mat3:
			gl.UniformMatrix3fv(loc, v)
		case []mgl32.Mat4:
			gl.UniformMatrix4fv(loc, v)
		default:
			panic(errors.New("invalid type"))
		}
		return true
	}
}

func compileShaderSource(source string, xtype gl.Enum) (gl.Shader, error) {
	id := gl.CreateShader(xtype)
	gl.ShaderSource(id, source)
	gl.CompileShader(id)

	if gl.GetShaderi(id, gl.COMPILE_STATUS) == gl.FALSE {
		return 0, fmt.Errorf(gl.GetShaderInfoLog(id))
	}

	return id, nil
}

func NewShaderProgram(vshader, fshader string) (*ShaderProgram, error) {
	if vref, err := compileShaderSource(vshader, gl.VERTEX_SHADER); err != nil {
		return nil, err
	} else if fref, err := compileShaderSource(fshader, gl.FRAGMENT_SHADER); err != nil {
		return nil, err
	} else {
		defer gl.DeleteShader(vref)
		defer gl.DeleteShader(fref)

		prog := &ShaderProgram{
			id:    gl.CreateProgram(),
			attrs: map[string]gl.Uniform{},
		}
		runtime.SetFinalizer(prog, (*ShaderProgram).delete)

		gl.AttachShader(prog.id, vref)
		gl.AttachShader(prog.id, fref)
		gl.LinkProgram(prog.id)

		if gl.GetProgrami(prog.id, gl.LINK_STATUS) == gl.FALSE {
			return nil, fmt.Errorf(gl.GetProgramInfoLog(prog.id))
		} else {
			return prog, nil
		}
	}
}
