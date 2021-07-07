package pancake2d

import (
	"image/color"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/mathx"
	gl "github.com/askeladdk/pancake/opengl"
)

var vertexFormat = pancake.AttribFormat{
	pancake.AttribVec2,    // XY
	pancake.AttribVec2,    // UV
	pancake.AttribByte4,   // RGBA
	pancake.AttribFloat32, // Z
}

type vertex struct {
	XY   mathx.Vec2
	UV   mathx.Vec2
	RGBA color.RGBA
	Z    float32
}

var (
	quadVertices = []vertex{
		{
			XY: mathx.Vec2{-.5, -.5},
			UV: mathx.Vec2{0, 0},
		},
		{
			XY: mathx.Vec2{-.5, +.5},
			UV: mathx.Vec2{0, 1},
		},
		{
			XY: mathx.Vec2{+.5, -.5},
			UV: mathx.Vec2{1, 0},
		},
		{
			XY: mathx.Vec2{+.5, +.5},
			UV: mathx.Vec2{1, 1},
		},
	}

	quadIndices = []uint32{0, 1, 2, 1, 2, 3}
)

// SpriteBatch is a set of sprites that can be drawn to a SpriteBuffer.
type SpriteBatch interface {
	Len() int
	TintColorAt(i int) color.Color
	TextureAt(i int) *pancake.Texture
	TextureRegionAt(i int) pancake.TextureRegion
	ModelViewAt(i int) mathx.Aff3
	OriginAt(i int) mathx.Vec2
	ZOrderAt(i int) float64
}

func toRGBA(c color.Color) color.RGBA {
	switch v := c.(type) {
	case color.RGBA:
		return v
	default:
		r, g, b, a := v.RGBA()
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
}

func instanceVertices(batch SpriteBatch, lo, hi int, vertices []vertex) []vertex {
	tmpv := make([]vertex, len(quadVertices))

	for i := lo; i < hi; i++ {
		modelview := batch.ModelViewAt(i)
		region := batch.TextureRegionAt(i).Aff3()
		rgba := toRGBA(batch.TintColorAt(i))
		origin := batch.OriginAt(i)
		z := float32(batch.ZOrderAt(i))
		for j, v := range quadVertices {
			tmpv[j] = vertex{
				XY:   modelview.Project(v.XY.Sub(origin)),
				UV:   region.Project(v.UV),
				RGBA: rgba,
				Z:    z,
			}
		}
		for _, n := range quadIndices {
			vertices = append(vertices, tmpv[n])
		}
	}

	return vertices
}

// SpriteDrawer renders SpriteBuffers.
type SpriteDrawer struct {
	vertices []vertex
	vbuffer  *pancake.VertexBuffer
	vslice   *pancake.VertexArraySlice
}

// NewSpriteDrawer creates a new sprite buffer with given capacity.
func NewSpriteDrawer(capacity int) *SpriteDrawer {
	vbuffer := pancake.NewVertexBuffer(vertexFormat, capacity*len(quadIndices), nil)
	vslice := pancake.NewVertexArraySlice(vbuffer)
	return &SpriteDrawer{
		vbuffer: vbuffer,
		vslice:  vslice,
	}
}

// Draw renders a SpriteBatch.
func (d *SpriteDrawer) Draw(batch SpriteBatch) {
	if batch.Len() == 0 {
		return
	}

	d.vslice.Begin()

	lo, texture := 0, batch.TextureAt(0)
	for hi := 1; hi < batch.Len(); hi++ {
		hitexture := batch.TextureAt(hi)
		if hitexture != texture {
			d.vertices = instanceVertices(batch, lo, hi, d.vertices[:0])
			texture.Begin()
			d.drawVertices(gl.TRIANGLES, d.vertices)
			texture.End()
			lo, texture = hi, hitexture
		}
	}

	hi := batch.Len()
	d.vertices = instanceVertices(batch, lo, hi, d.vertices[:0])
	texture.Begin()
	d.drawVertices(gl.TRIANGLES, d.vertices)
	texture.End()

	d.vslice.End()
}

func (d *SpriteDrawer) drawVertices(mode gl.Enum, verts []vertex) {
	lo, step := 0, d.vslice.Len()
	for hi := step; hi < len(verts); hi += step {
		vslice := d.vslice.Slice(0, hi-lo)
		vslice.SetData(verts[lo:hi])
		vslice.Draw(mode)
		lo = hi
	}

	hi := len(verts)
	vslice := d.vslice.Slice(0, hi-lo)
	vslice.SetData(verts[lo:hi])
	vslice.Draw(mode)
}
