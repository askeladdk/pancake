package text

import (
	"fmt"
	"testing"

	"golang.org/x/image/font/basicfont"
)

func TestNewText(t *testing.T) {
	font := NewFont(basicfont.Face7x13, ASCII)
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
