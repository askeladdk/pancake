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

var c mgl.Context3

func Init(param interface{}) {
	if c3, ok := param.(gl.Context3); !ok {
		panic(errors.New("gl package not compiled with ES3 support"))
	} else {
		c = c3
	}
}
