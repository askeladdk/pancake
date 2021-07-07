package pancake2d

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/askeladdk/pancake"
	"golang.org/x/image/font/basicfont"
)

func TestMain(m *testing.M) {
	pancake.Main(pancake.Options{
		WindowSize: image.Point{320, 200},
	}, func(_ pancake.App) error {
		os.Exit(m.Run())
		return nil
	})
}

func TestNewText(t *testing.T) {
	font := pancake.NewFont(basicfont.Face7x13, pancake.ASCII)
	text := NewText(font)
	fmt.Fprintf(text, "hello\r\n\t")
	text.WriteString("world")
	if text.Len() != 10 {
		t.Fatal()
	}

	_ = text.ModelViewAt(0)
	_ = text.OriginAt(0)
	_ = text.TextureAt(0)
	_ = text.TextureRegionAt(0)
	_ = text.TintColorAt(0)
	_ = text.ZOrderAt(0)

	text.Reset()
	if text.Len() != 0 {
		t.Fatal()
	}
}
