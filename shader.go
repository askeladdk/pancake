package pancake

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/askeladdk/pancake/mathx"
	gl "github.com/askeladdk/pancake/opengl"
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
		case float64:
			gl.Uniform1f(loc, v)
		case mathx.Vec2:
			gl.Uniform2fv(loc, []mathx.Vec2{v})
		case mathx.Vec3:
			gl.Uniform3fv(loc, []mathx.Vec3{v})
		case mathx.Vec4:
			gl.Uniform4fv(loc, []mathx.Vec4{v})
		case mathx.Mat3:
			gl.UniformMatrix3fv(loc, []mathx.Mat3{v})
		case mathx.Mat4:
			gl.UniformMatrix4fv(loc, []mathx.Mat4{v})
		case mathx.Aff3:
			gl.UniformMatrix3fv(loc, []mathx.Mat3{v.Mat3()})
		case []float64:
			gl.Uniform1fv(loc, v)
		case []mathx.Vec2:
			gl.Uniform2fv(loc, v)
		case []mathx.Vec3:
			gl.Uniform3fv(loc, v)
		case []mathx.Vec4:
			gl.Uniform4fv(loc, v)
		case []mathx.Mat3:
			gl.UniformMatrix3fv(loc, v)
		case []mathx.Mat4:
			gl.UniformMatrix4fv(loc, v)
		case []mathx.Aff3:
			xs := make([]mathx.Mat3, len(v))
			for i, x := range v {
				xs[i] = x.Mat3()
			}
			gl.UniformMatrix3fv(loc, xs)
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
