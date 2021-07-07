package pancake

import (
	"image"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Main(Options{
		WindowSize: image.Point{320, 200},
	}, func(_ App) error {
		os.Exit(m.Run())
		return nil
	})
}
