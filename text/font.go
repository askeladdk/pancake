package text

import (
	"image"
	"image/draw"
	"unicode"

	"github.com/askeladdk/pancake/mathx"

	"github.com/askeladdk/pancake/graphics"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// ASCII is the set of all printable ASCII characters and the replacement character.
var ASCII = &unicode.RangeTable{
	LatinOffset: 1,
	R16: []unicode.Range16{
		{Lo: ' ', Hi: '~', Stride: 1},
		{Lo: 0xFFFD, Hi: 0xFFFD, Stride: 1},
	},
}

func fixedToFloat64(x fixed.Int26_6) float64 {
	return float64(x) / (1 << 6)
}

func rangeTableToRunes(rangeTab *unicode.RangeTable) []rune {
	var runes []rune
	for _, r16 := range rangeTab.R16 {
		for c := r16.Lo; c <= r16.Hi; c += r16.Stride {
			runes = append(runes, rune(c))
		}
	}
	for _, r32 := range rangeTab.R32 {
		for c := r32.Lo; c <= r32.Hi; c += r32.Stride {
			runes = append(runes, rune(c))
		}
	}
	return runes
}

// Glyph represents a character in a Face.
type Glyph struct {
	// Region is the rectangular area of the font texture that contains the glyph.
	Region graphics.TextureRegion

	// Scale is the size of the character.
	Scale mathx.Vec2

	// Advance is the number of pixels to advance horizontally to the character.
	Advance float64
}

// Font is a renderable font.Font.
type Font struct {
	face       font.Face
	texture    *graphics.Texture
	mapping    map[rune]Glyph
	lineHeight float64
}

// NewFont creates a new Font from a font.Face by building a texture atlas
// containing all characters in the range table.
func NewFont(face font.Face, rangeTab *unicode.RangeTable) *Font {
	runes := rangeTableToRunes(rangeTab)

	metrics := face.Metrics()
	padding := fixed.I(2)
	height := metrics.Ascent + metrics.Descent
	width := fixed.I(0)
	for _, r := range runes {
		if b, _, ok := face.GlyphBounds(r); ok {
			width += fixed.I((b.Max.X - b.Min.X).Ceil())
			width += padding
		}
	}

	imageSize := image.Point{width.Ceil(), height.Ceil()}

	rgba := image.NewRGBA(image.Rectangle{image.Point{}, imageSize})

	mapping := map[rune]Glyph{}

	dot := fixed.Point26_6{
		X: 0,
		Y: face.Metrics().Ascent,
	}

	for _, r := range runes {
		if b, a, ok := face.GlyphBounds(r); ok {
			w := (b.Max.X - b.Min.X).Ceil() + padding.Ceil()
			h := imageSize.Y
			x := dot.X.Ceil()
			region := image.Rectangle{
				image.Point{x, 0},
				image.Point{x + w, h},
			}

			mapping[r] = Glyph{
				Region:  graphics.NewTextureRegion(imageSize, region),
				Scale:   mathx.Vec2{float64(w), float64(h)},
				Advance: fixedToFloat64(a),
			}

			dr, mask, maskp, _, _ := face.Glyph(dot, r)
			draw.Draw(rgba, dr, mask, maskp, draw.Src)

			dot.X += fixed.I((b.Max.X - b.Min.X).Ceil())
			dot.X += padding
		}
	}

	return &Font{
		face:       face,
		mapping:    mapping,
		texture:    graphics.NewTextureFromImage(rgba, graphics.FilterLinear),
		lineHeight: float64(face.Metrics().Height.Ceil()),
	}
}

// Face returns the source font.Face.
func (fnt *Font) Face() font.Face {
	return fnt.face
}

// Texture returns the texture atlas that contains all glyphs.
func (fnt *Font) Texture() *graphics.Texture {
	return fnt.texture
}

// Glyph returns the glyph associated with a rune.
func (fnt *Font) Glyph(r rune) Glyph {
	if glyph, ok := fnt.mapping[r]; ok {
		return glyph
	}
	return fnt.mapping[unicode.ReplacementChar]
}

// LineHeight reports the font line height in pixels.
func (fnt *Font) LineHeight() float64 {
	return fnt.lineHeight
}

// Kern reports the kerning distance in pixels between two runes.
func (fnt *Font) Kern(r0, r1 rune) float64 {
	return float64(fnt.face.Kern(r0, r1).Ceil())
}
