package text

import (
	"image/color"
	"unicode/utf8"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// Text represents the drawing state of text written in a Font.
// It implements graphics2d.Batch.
type Text struct {
	// Pos is the base position of the text.
	Pos mathx.Vec2

	// Dot is the current position of the cursor.
	Dot mathx.Vec2

	// TintColor is the font color.
	TintColor color.Color

	// Scale scales the entire text.
	Scale float64

	// TabWidth is the the space in pixels of the tab character.
	// Defaults to four times the width of the space character.
	TabWidth float64

	// ZOrder is the Z order.
	ZOrder float64

	modelview []mathx.Aff3
	region    []graphics.TextureRegion
	font      *Font
	lastRune  rune
}

// NewText creates a new Text batch based on a Font.
func NewText(font *Font) *Text {
	return &Text{
		TintColor: color.RGBA{255, 255, 255, 255},
		TabWidth:  font.Glyph(' ').Advance * 4,
		Scale:     1,
		font:      font,
	}
}

// Reset clears the Text.
func (t *Text) Reset() {
	t.Dot = mathx.Vec2{}
	t.modelview = t.modelview[:0]
	t.region = t.region[:0]
	t.lastRune = 0
}

// Write implements io.Writer.
func (t *Text) Write(p []byte) (int, error) {
	t.draw(p)
	return len(p), nil
}

// WriteString implements io.StringWriter.
func (t *Text) WriteString(s string) (int, error) {
	t.Write([]byte(s))
	return len(s), nil
}

func (t *Text) draw(buf []byte) {
	for utf8.FullRune(buf) {
		r, size := utf8.DecodeRune(buf)
		buf = buf[size:]

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

		advance := glyph.Advance
		if t.lastRune > 0 {
			advance += t.font.Kern(t.lastRune, r)
		}

		t.Dot[0] += advance * t.Scale

		t.lastRune = r
	}
}

// Len implements graphics2d.Batch.
func (t *Text) Len() int {
	return len(t.modelview)
}

// TintColorAt implements graphics2d.Batch.
func (t *Text) TintColorAt(i int) color.Color {
	return t.TintColor
}

// TextureAt implements graphics2d.Batch.
func (t *Text) TextureAt(_ int) *graphics.Texture {
	return t.font.Texture()
}

// TextureRegionAt implements graphics2d.Batch.
func (t *Text) TextureRegionAt(i int) graphics.TextureRegion {
	return t.region[i]
}

// ModelViewAt implements graphics2d.Batch.
func (t *Text) ModelViewAt(i int) mathx.Aff3 {
	return t.modelview[i].Translated(t.Pos)
}

// OriginAt implements graphics2d.Batch.
func (t *Text) OriginAt(i int) mathx.Vec2 {
	return mathx.Vec2{-.5, -.5} // top left corner
}

// ZOrderAt implements graphics2d.Batch.
func (t *Text) ZOrderAt(i int) float64 {
	return t.ZOrder
}
