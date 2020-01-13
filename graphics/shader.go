package graphics

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/faiface/mainthread"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var shaderBinder = newBinder(func(prog uint32) {
	gl.UseProgram(prog)
})

type shader struct {
	ref   uint32
	attrs map[string]int32
}

func (this *shader) Begin() {
	shaderBinder.bind(this.ref)
}

func (this *shader) End() {
	shaderBinder.unbind()
}

func (this *shader) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteProgram(this.ref)
	})
}

func (this *shader) getUniformLocation(name string) int32 {
	if loc, ok := this.attrs[name]; !ok {
		loc = gl.GetUniformLocation(this.ref, gl.Str(name+"\x00"))
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
		case int32:
			gl.Uniform1i(loc, v)
		case uint32:
			gl.Uniform1ui(loc, v)
		case float32:
			gl.Uniform1f(loc, v)
		case mgl32.Vec2:
			gl.Uniform2fv(loc, 1, &v[0])
		case mgl32.Vec3:
			gl.Uniform3fv(loc, 1, &v[0])
		case mgl32.Vec4:
			gl.Uniform4fv(loc, 1, &v[0])
		case mgl32.Mat2:
			gl.UniformMatrix2fv(loc, 1, false, &v[0])
		case mgl32.Mat3:
			gl.UniformMatrix3fv(loc, 1, false, &v[0])
		case mgl32.Mat4:
			gl.UniformMatrix4fv(loc, 1, false, &v[0])
		default:
			panic(errors.New("invalid type"))
		}
		return true
	}
}

func compileShaderSource(source string, xtype uint32) (uint32, error) {
	ref := gl.CreateShader(xtype)
	src, free := gl.Strs(source)
	srclen := int32(len(source))
	defer free()
	gl.ShaderSource(ref, 1, src, &srclen)
	gl.CompileShader(ref)

	var success int32
	gl.GetShaderiv(ref, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(ref, gl.INFO_LOG_LENGTH, &logLen)
		infoLog := make([]byte, logLen)
		gl.GetShaderInfoLog(ref, logLen, nil, &infoLog[0])
		return 0, fmt.Errorf(string(infoLog))
	}

	return ref, nil
}

func newShader(vshader, fshader string) (*shader, error) {
	if vref, err := compileShaderSource(vshader, gl.VERTEX_SHADER); err != nil {
		return nil, err
	} else if fref, err := compileShaderSource(fshader, gl.FRAGMENT_SHADER); err != nil {
		return nil, err
	} else {
		defer gl.DeleteShader(vref)
		defer gl.DeleteShader(fref)

		prog := &shader{
			ref:   gl.CreateProgram(),
			attrs: map[string]int32{},
		}
		runtime.SetFinalizer(prog, (*shader).delete)

		gl.AttachShader(prog.ref, vref)
		gl.AttachShader(prog.ref, fref)
		gl.LinkProgram(prog.ref)

		var success int32
		gl.GetProgramiv(prog.ref, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(prog.ref, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(prog.ref, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf(string(infoLog))
		} else {
			return prog, nil
		}
	}
}
