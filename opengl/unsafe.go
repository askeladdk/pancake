package opengl

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
//  var data []uint8
// 	...
// 	gl.TexImage2D(gl.TEXTURE_2D, ..., gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
//
// The implementation is taken from
// https://github.com/go-gl
func Ptr(data interface{}) unsafe.Pointer {
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)
	if data == nil || v.IsNil() {
		return unsafe.Pointer(nil)
	}
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		default:
			panic(fmt.Errorf("unsupported pointer to type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", e.Kind()))
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(v.Pointer())
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	default:
		panic(fmt.Errorf("unsupported type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", v.Type()))
	}
	return addr
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
//
// The implementation is taken from
// https://github.com/go-gl
func ptrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}

// PtrToSlice creates a byte slice from an arbitrary pointer.
// Useful for functions in golang.org/x/mobile/gl that take
// a []byte for what should be an unsafe.Pointer.
//
// The implementation is taken from
// https://github.com/golang/sys/blob/86b910548bc16777f40503131aa424ae0a092199/unix/syscall_unix.go#L116-L124
func ptrToSlice(ptr unsafe.Pointer, length int) []byte {
	sl := reflect.SliceHeader{uintptr(ptr), length, length}
	return *(*[]byte)(unsafe.Pointer(&sl))
}
