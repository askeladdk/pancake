package text

import (
	"image/color"
	"unicode/utf8"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

type Text struct {
	Pos        mathx.Vec2
	Dot        mathx.Vec2
	Color      color.NRGBA
	Scale      float32
	LineHeight float32
	TabWidth   float32

	modelview []mathx.Aff3
	region    []graphics.TextureRegion
	colors    []color.NRGBA
	buffer    []byte
	font      Font
	lastRune  rune
}

func NewText(font Font) *Text {
	return &Text{
		Color:      color.NRGBA{255, 255, 255, 255},
		LineHeight: font.LineHeight(),
		TabWidth:   font.Glyph(' ').Advance * 4,
		Scale:      1,
		font:       font,
		lastRune:   -1,
	}
}

func (t *Text) Clear() {
	t.Dot = mathx.Vec2{}
	t.modelview = t.modelview[:0]
	t.region = t.region[:0]
	t.colors = t.colors[:0]
	t.buffer = t.buffer[:0]
	t.lastRune = -1
}

func (t *Text) Write(p []byte) (int, error) {
	t.buffer = append(t.buffer, p...)
	t.draw()
	return len(p), nil
}

func (t *Text) WriteString(s string) (int, error) {
	t.buffer = append(t.buffer, s...)
	t.draw()
	return len(s), nil
}

func (t *Text) draw() {
	for utf8.FullRune(t.buffer) {
		r, size := utf8.DecodeRune(t.buffer)
		t.buffer = t.buffer[size:]

		switch r {
		case '\n':
			t.Dot[0] = 0
			t.Dot[1] += t.font.LineHeight() * t.Scale
			continue
		case '\r':
			t.Dot[0] = 0
			continue
		case '\t':
			t.Dot[0] += t.TabWidth * t.Scale
			continue
		}

		glyph := t.font.Glyph(r)

		t.modelview = append(t.modelview, mathx.
			ScaleAff3(glyph.Scale.Mul(t.Scale)).
			Translated(t.Dot),
		)

		t.region = append(t.region, glyph.Region)
		t.colors = append(t.colors, t.Color)

		advance := glyph.Advance
		if t.lastRune >= 0 {
			advance += t.font.Kern(t.lastRune, r)
		}

		t.Dot[0] += advance * t.Scale

		t.lastRune = r
	}
}

func (t *Text) Len() int {
	return len(t.modelview)
}

func (t *Text) TintColorAt(i int) color.NRGBA {
	return t.colors[i]
}

func (t *Text) Texture() *graphics.Texture {
	return t.font.Texture()
}

func (t *Text) TextureRegionAt(i int) graphics.TextureRegion {
	return t.region[i]
}

func (t *Text) ModelViewAt(i int) mathx.Aff3 {
	return t.modelview[i].Translated(t.Pos)
}

func (t *Text) OriginAt(i int) mathx.Vec2 {
	return mathx.Vec2{-.5, -.5}
}
