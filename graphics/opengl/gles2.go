// +build android ios
// +build !darwin
// +build !freebsd
// +build !linux
// +build !windows

package opengl

import (
	"errors"

	"golang.org/x/mobile/gl"
)

var c mgl.Context

func Init(param interface{}) {
	if c3, ok := param.(gl.Context); !ok {
		panic(errors.New("gl package not compiled with ES2 support"))
	} else {
		c = c3
	}
}
