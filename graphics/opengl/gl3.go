// +build darwin freebsd linux windows
// +build !android
// +build !ios

package opengl

import (
	"image"
	"unsafe"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func Init(param interface{}) error {
	return gl.Init()
}

func BindBuffer(target Enum, buffer Buffer) {
	mainthread.Call(func() {
		gl.BindBuffer(uint32(target), uint32(buffer))
	})
}

func CreateBuffer() Buffer {
	var buffer uint32
	mainthread.Call(func() {
		gl.GenBuffers(1, &buffer)
	})
	return Buffer(buffer)
}

func DeleteBuffer(buffer Buffer) {
	mainthread.CallNonBlock(func() {
		gl.DeleteBuffers(1, (*uint32)(&buffer))
	})
}

func BindFramebuffer(target Enum, frame Framebuffer) {
	mainthread.Call(func() {
		gl.BindFramebuffer(uint32(target), uint32(frame))
	})
}

func CreateFramebuffer() Framebuffer {
	var frame uint32
	mainthread.Call(func() {
		gl.GenFramebuffers(1, &frame)
	})
	return Framebuffer(frame)
}

func DeleteFramebuffer(frame Framebuffer) {
	mainthread.CallNonBlock(func() {
		gl.DeleteFramebuffers(1, (*uint32)(&frame))

	})
}

func BindRenderbuffer(rbuf Renderbuffer) {
	mainthread.Call(func() {
		gl.BindRenderbuffer(gl.RENDERBUFFER, uint32(rbuf))
	})
}

func CreateRenderbuffer() Renderbuffer {
	var rbuf uint32
	mainthread.Call(func() {
		gl.GenRenderbuffers(1, &rbuf)
	})
	return Renderbuffer(rbuf)
}

func DeleteRenderbuffer(rbuf Renderbuffer) {
	mainthread.CallNonBlock(func() {
		gl.GenRenderbuffers(1, (*uint32)(&rbuf))
	})
}

func BindTexture(target Enum, texture Texture) {
	mainthread.Call(func() {
		gl.BindTexture(uint32(target), uint32(texture))
	})
}

func CreateTexture() Texture {
	var texture uint32
	mainthread.Call(func() {
		gl.GenTextures(1, &texture)
	})
	return Texture(texture)
}

func DeleteTexture(texture Texture) {
	mainthread.CallNonBlock(func() {
		gl.DeleteTextures(1, (*uint32)(&texture))
	})
}

func BindVertexArray(array VertexArray) {
	mainthread.Call(func() {
		gl.BindVertexArray(uint32(array))
	})
}

func CreateVertexArray() VertexArray {
	var array uint32
	mainthread.Call(func() {
		gl.GenVertexArrays(1, &array)
	})
	return VertexArray(array)
}

func DeleteVertexArray(array VertexArray) {
	mainthread.CallNonBlock(func() {
		gl.DeleteVertexArrays(1, (*uint32)(&array))
	})
}

func BindProgram(program Program) {
	mainthread.Call(func() {
		gl.UseProgram(uint32(program))
	})
}

func CreateProgram() Program {
	var program uint32
	mainthread.Call(func() {
		program = gl.CreateProgram()
	})
	return Program(program)
}

func DeleteProgram(program Program) {
	mainthread.CallNonBlock(func() {
		gl.DeleteProgram(uint32(program))
	})
}

func CreateShader(xtype Enum) Shader {
	var shader uint32
	mainthread.Call(func() {
		shader = gl.CreateShader(uint32(xtype))
	})
	return Shader(shader)
}

func DeleteShader(shader Shader) {
	mainthread.CallNonBlock(func() {
		gl.DeleteShader(uint32(shader))
	})
}

func Enable(cap Enum) {
	mainthread.Call(func() {
		gl.Enable(uint32(cap))
	})
}

func Disable(cap Enum) {
	mainthread.Call(func() {
		gl.Disable(uint32(cap))
	})
}

func BlendFunc(sfactor, dfactor Enum) {
	mainthread.Call(func() {
		gl.BlendFunc(uint32(sfactor), uint32(dfactor))
	})
}

func GetInteger(name Enum) int {
	var data int32
	mainthread.Call(func() {
		gl.GetIntegerv(uint32(name), &data)
	})
	return int(data)
}

func GetString(name Enum) string {
	var str string
	mainthread.Call(func() {
		str = gl.GoStr(gl.GetString(uint32(name)))
	})
	return str
}

func Viewport(r image.Rectangle) {
	mainthread.Call(func() {
		size := r.Size()
		gl.Viewport(int32(r.Min.X), int32(r.Min.Y), int32(size.X), int32(size.Y))
	})
}

func Clear(mask Enum) {
	mainthread.Call(func() {
		gl.Clear(uint32(mask))
	})
}

func GetError() Enum {
	var err uint32
	mainthread.Call(func() {
		err = gl.GetError()
	})
	return Enum(err)
}

func BlitNamedFramebuffer(src, dst Framebuffer, sr, dr image.Rectangle, mask, filter Enum) {
	mainthread.Call(func() {
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, uint32(dst))
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, uint32(src))
		gl.BlitFramebuffer(
			int32(sr.Min.X), int32(sr.Min.Y), int32(sr.Max.X), int32(sr.Max.Y),
			int32(dr.Min.X), int32(dr.Min.Y), int32(dr.Max.X), int32(dr.Max.Y),
			uint32(mask), uint32(filter),
		)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	})
}

func FramebufferTexture2D(attachment, textarget Enum, texture Texture, level int) {
	mainthread.Call(func() {
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, uint32(attachment),
			uint32(textarget), uint32(texture), int32(level))
	})
}

func FramebufferRenderbuffer(attachment, rbuffer Renderbuffer) {
	mainthread.Call(func() {
		gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, uint32(attachment),
			gl.RENDERBUFFER, uint32(rbuffer))
	})
}

func CheckFramebufferStatus() Enum {
	var status uint32
	mainthread.Call(func() {
		status = gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	})
	return Enum(status)
}

func GetUniformLocation(program Program, name string) Uniform {
	var loc int32
	mainthread.Call(func() {
		loc = gl.GetUniformLocation(uint32(program), gl.Str(name+"\x00"))
	})
	return Uniform(loc)
}

func Uniform1i(dst Uniform, v0 int) {
	mainthread.Call(func() {
		gl.Uniform1i(int32(dst), int32(v0))
	})
}

func Uniform1ui(dst Uniform, v0 uint) {
	mainthread.Call(func() {
		gl.Uniform1ui(int32(dst), uint32(v0))
	})
}

func Uniform1f(dst Uniform, v0 float32) {
	mainthread.Call(func() {
		gl.Uniform1f(int32(dst), v0)
	})
}

func Uniform1fv(dst Uniform, v []float32) {
	mainthread.Call(func() {
		gl.Uniform1fv(int32(dst), int32(len(v)), &v[0])
	})
}

func Uniform2fv(dst Uniform, vs []mgl32.Vec2) {
	mainthread.Call(func() {
		gl.Uniform2fv(int32(dst), int32(2*len(vs)), &vs[0][0])
	})
}

func Uniform3fv(dst Uniform, vs []mgl32.Vec3) {
	mainthread.Call(func() {
		gl.Uniform3fv(int32(dst), int32(3*len(vs)), &vs[0][0])
	})
}

func Uniform4fv(dst Uniform, vs []mgl32.Vec4) {
	mainthread.Call(func() {
		gl.Uniform4fv(int32(dst), int32(4*len(vs)), &vs[0][0])
	})
}

func UniformMatrix2fv(dst Uniform, vs []mgl32.Mat2) {
	mainthread.Call(func() {
		gl.UniformMatrix2fv(int32(dst), int32(len(vs)), false, &vs[0][0])
	})
}

func UniformMatrix3fv(dst Uniform, vs []mgl32.Mat3) {
	mainthread.Call(func() {
		gl.UniformMatrix3fv(int32(dst), int32(len(vs)), false, &vs[0][0])
	})
}

func UniformMatrix4fv(dst Uniform, vs []mgl32.Mat4) {
	mainthread.Call(func() {
		gl.UniformMatrix4fv(int32(dst), int32(len(vs)), false, &vs[0][0])
	})
}

func ShaderSource(shader Shader, source string) {
	mainthread.Call(func() {
		src, free := gl.Strs(source)
		srclen := int32(len(source))
		defer free()
		gl.ShaderSource(uint32(shader), 1, src, &srclen)
	})
}

func CompileShader(shader Shader) {
	mainthread.Call(func() {
		gl.CompileShader(uint32(shader))
	})
}

func GetShaderi(shader Shader, pname Enum) int {
	var v int32
	mainthread.Call(func() {
		gl.GetShaderiv(uint32(shader), uint32(pname), &v)
	})
	return int(v)
}

func GetShaderInfoLog(shader Shader) string {
	var str string
	mainthread.Call(func() {
		var logLen int32
		gl.GetShaderiv(uint32(shader), gl.INFO_LOG_LENGTH, &logLen)
		infoLog := make([]byte, logLen)
		gl.GetShaderInfoLog(uint32(shader), logLen, nil, &infoLog[0])
		str = string(infoLog)
	})
	return str
}

func AttachShader(program Program, shader Shader) {
	mainthread.Call(func() {
		gl.AttachShader(uint32(program), uint32(shader))
	})
}

func LinkProgram(program Program) {
	mainthread.Call(func() {
		gl.LinkProgram(uint32(program))
	})
}

func GetProgrami(program Program, pname Enum) int {
	var v int32
	mainthread.Call(func() {
		gl.GetProgramiv(uint32(program), uint32(pname), &v)
	})
	return int(v)
}

func GetProgramInfoLog(program Program) string {
	var str string
	mainthread.Call(func() {
		var logLen int32
		gl.GetProgramiv(uint32(program), gl.INFO_LOG_LENGTH, &logLen)
		infoLog := make([]byte, logLen)
		gl.GetProgramInfoLog(uint32(program), logLen, nil, &infoLog[0])
		str = string(infoLog)
	})
	return str
}

func ActiveTexture(target Enum) {
	mainthread.Call(func() {
		gl.ActiveTexture(uint32(target))
	})
}

func GenerateMipmap(target Enum) {
	mainthread.Call(func() {
		gl.GenerateMipmap(uint32(target))
	})
}

func TexParameteri(target, pname, param Enum) {
	mainthread.Call(func() {
		gl.TexParameteri(uint32(target), uint32(pname), int32(param))
	})
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, xtype Enum, data []byte) {
	mainthread.Call(func() {
		gl.TexSubImage2D(uint32(target), int32(level),
			int32(x), int32(y), int32(width), int32(height),
			uint32(format), uint32(xtype), Ptr(data))
	})
}

func GetTexImage(target Enum, level int, format, xtype Enum, data []byte) {
	mainthread.Call(func() {
		gl.GetTexImage(uint32(target), int32(level), uint32(format), uint32(xtype), Ptr(data))
	})
}

func TexImage2D(target Enum, level int, internalFormat Enum, width, height int, format, xtype Enum, data []byte) {
	mainthread.Call(func() {
		gl.TexImage2D(
			uint32(target),
			int32(level),
			int32(internalFormat),
			int32(width),
			int32(height),
			0,
			uint32(format),
			uint32(xtype),
			Ptr(data),
		)
	})
}

func VertexAttribPointer(dst Attrib, size int, xtype Enum, normalized bool, stride, offset int) {
	mainthread.Call(func() {
		gl.VertexAttribPointer(uint32(dst), int32(size), uint32(xtype),
			normalized, int32(stride), ptrOffset(offset))
	})
}

func VertexAttribDivisor(dst Attrib, divisor int) {
	mainthread.Call(func() {
		gl.VertexAttribDivisor(uint32(dst), uint32(divisor))
	})
}

func EnableVertexAttribArray(dst Attrib) {
	mainthread.Call(func() {
		gl.EnableVertexAttribArray(uint32(dst))
	})
}

func DrawArrays(mode Enum, first, count int) {
	mainthread.Call(func() {
		gl.DrawArrays(uint32(mode), int32(first), int32(count))
	})
}

func DrawArraysInstanced(mode Enum, first, count, instances int) {
	mainthread.Call(func() {
		gl.DrawArraysInstanced(uint32(mode), int32(first), int32(count), int32(instances))
	})
}

func DrawElements(mode Enum, count int, xtype Enum, indices int) {
	mainthread.Call(func() {
		gl.DrawElements(uint32(mode), int32(count), uint32(xtype), ptrOffset(indices))
	})
}

func DrawElementsInstanced(mode Enum, count int, xtype Enum, indices, instanceCount int) {
	mainthread.Call(func() {
		gl.DrawElementsInstanced(uint32(mode), int32(count), uint32(xtype), ptrOffset(indices), int32(instanceCount))
	})
}

func BufferSubData(target Enum, offset, size int, data unsafe.Pointer) {
	mainthread.Call(func() {
		gl.BufferSubData(uint32(target), offset, size, data)
	})
}

func BufferData(target Enum, len int, data unsafe.Pointer, usage Enum) {
	mainthread.Call(func() {
		gl.BufferData(uint32(target), len, data, uint32(usage))
	})
}

func Flush() {
	mainthread.Call(gl.Flush)
}

func ClearColor(r, g, b, a float32) {
	mainthread.Call(func() {
		gl.ClearColor(r, g, b, a)
	})
}

func RenderbufferStorage(internalFormat Enum, width, height, samples int) {
	mainthread.Call(func() {
		gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, int32(samples), uint32(internalFormat), int32(width), int32(height))
	})
}

func Scissor(r image.Rectangle) {
	mainthread.Call(func() {
		size := r.Size()
		gl.Scissor(int32(r.Min.X), int32(r.Min.Y), int32(size.X), int32(size.Y))
	})
}
