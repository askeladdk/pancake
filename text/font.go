package text

import (
	"image"
	"image/draw"
	"sort"
	"unicode"

	"github.com/askeladdk/pancake/mathx"

	"github.com/askeladdk/pancake/graphics"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var ASCII = []rune(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")

func fixedToFloat32(x fixed.Int26_6) float32 {
	return float32(x) / (1 << 6)
}

func combineRuneSets(runeSets [][]rune) []rune {
	seen := map[rune]struct{}{}
	runes := []rune{unicode.ReplacementChar}
	for _, set := range runeSets {
		for _, r := range set {
			if _, ok := seen[r]; !ok {
				runes = append(runes, r)
				seen[r] = struct{}{}
			}
		}
	}

	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return runes
}

type Font interface {
	Texture() *graphics.Texture
	Glyph(r rune) Glyph
	LineHeight() float32
	Kern(r0, r1 rune) float32
}

type Glyph struct {
	Region  graphics.TextureRegion
	Scale   mathx.Vec2
	Advance float32
}

type faceFont struct {
	face       font.Face
	texture    *graphics.Texture
	mapping    map[rune]Glyph
	lineHeight float32
}

func NewFontFromFace(face font.Face, runeSets ...[]rune) Font {
	runes := combineRuneSets(runeSets)

	metrics := face.Metrics()
	padding := fixed.I(2)
	height := metrics.Ascent + metrics.Descent
	width := fixed.I(0)
	for _, r := range runes {
		if b, _, ok := face.GlyphBounds(r); ok {
			width += fixed.I((b.Max.X - b.Min.X).Ceil())
			// width = width.Ceil())
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
				Region: graphics.NewTextureRegion(imageSize, region),
				Scale: mathx.Vec2{
					float32(w),
					float32(h),
				},
				Advance: fixedToFloat32(a),
			}

			dr, mask, maskp, _, _ := face.Glyph(dot, r)
			draw.Draw(rgba, dr, mask, maskp, draw.Src)

			dot.X += fixed.I((b.Max.X - b.Min.X).Ceil())
			// dot.X = fixed.I(dot.X.Ceil())
			dot.X += padding
		}
	}

	return &faceFont{
		face:       face,
		mapping:    mapping,
		texture:    graphics.NewTextureFromImage(rgba, graphics.FilterLinear),
		lineHeight: float32(face.Metrics().Height.Ceil()),
	}
}

func (ttf *faceFont) Texture() *graphics.Texture {
	return ttf.texture
}

func (ttf *faceFont) Glyph(r rune) Glyph {
	if glyph, ok := ttf.mapping[r]; ok {
		return glyph
	} else if glyph, ok := ttf.mapping[unicode.ReplacementChar]; ok {
		return glyph
	} else {
		return Glyph{}
	}
}

func (ttf *faceFont) LineHeight() float32 {
	return ttf.lineHeight
}

func (ttf *faceFont) Kern(r0, r1 rune) float32 {
	return float32(ttf.face.Kern(r0, r1).Ceil())
}
