package pancake

import (
	"testing"
	"unicode"

	"golang.org/x/image/font/basicfont"
)

func Test_rangeTableToRunes(t *testing.T) {
	expected := []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\ufffd")
	runes := rangeTableToRunes(ASCII)
	if len(runes) != len(expected) {
		t.Fatal()
	}
	for i, c := range expected {
		if runes[i] != c {
			t.Fatal(c)
		}
	}

	var asciir32 = &unicode.RangeTable{
		R32: []unicode.Range32{
			{Lo: ' ', Hi: '~', Stride: 1},
			{Lo: 0xFFFD, Hi: 0xFFFD, Stride: 1},
		},
	}

	runes = rangeTableToRunes(asciir32)
	for i, c := range expected {
		if runes[i] != c {
			t.Fatal(c)
		}
	}
}

func TestNewFont(t *testing.T) {
	font := NewFont(basicfont.Face7x13, ASCII)
	if font.LineHeight() != 13 {
		t.Fatal()
	} else if len(font.mapping) != 96 {
		t.Fatal()
	}
	_ = font.Face()
	_ = font.Texture()
	_ = font.Kern('a', 'b')
	_ = font.Glyph('a')
	_ = font.Glyph('€')
}
