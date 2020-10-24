package graphics2d

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/mathx"
)

var VertexFormat = graphics.AttribFormat{
	graphics.Vec2,    // XY
	graphics.Vec2,    // UV
	graphics.Byte4,   // RGBA
	graphics.Float32, // Z
}

type Vertex struct {
	XY   mathx.Vec2
	UV   mathx.Vec2
	RGBA color.RGBA
	Z    float32
}

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	DrawMode gl.Enum
}

type Batch interface {
	Len() int
	TintColorAt(i int) color.Color
	TextureAt(i int) *graphics.Texture
	TextureRegionAt(i int) graphics.TextureRegion
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

func instanceVertices(mesh *Mesh, batch Batch, lo, hi int, vertices []Vertex) []Vertex {
	tmpv := make([]Vertex, len(mesh.Vertices))

	for i := lo; i < hi; i++ {
		modelview := batch.ModelViewAt(i)
		region := batch.TextureRegionAt(i).Aff3()
		rgba := toRGBA(batch.TintColorAt(i))
		origin := batch.OriginAt(i)
		z := float32(batch.ZOrderAt(i))
		for m, v := range mesh.Vertices {
			tmpv[m] = Vertex{
				XY:   modelview.Project(v.XY.Sub(origin)),
				UV:   region.Project(v.UV),
				RGBA: rgba,
				Z:    z,
			}
		}
		for _, n := range mesh.Indices {
			vertices = append(vertices, tmpv[n])
		}
	}

	return vertices
}

type Drawer struct {
	vertices []Vertex
	vbuffer  *graphics.Buffer
	vslice   *graphics.VertexSlice
	mesh     *Mesh
}

func NewDrawer(maxInstances int, mesh *Mesh) *Drawer {
	if mesh == nil {
		mesh = &Quad
	}

	vbuffer := graphics.NewBuffer(VertexFormat, maxInstances*len(mesh.Indices), nil)
	vslice := graphics.NewVertexSlice(vbuffer)
	return &Drawer{
		vbuffer: vbuffer,
		vslice:  vslice,
		mesh:    mesh,
	}
}

func (d *Drawer) Draw(batch Batch) {
	if batch.Len() == 0 {
		return
	}

	d.vslice.Begin()

	lo, texture := 0, batch.TextureAt(0)
	for hi := 1; hi < batch.Len(); hi++ {
		hitexture := batch.TextureAt(hi)
		if hitexture != texture {
			d.vertices = instanceVertices(d.mesh, batch, lo, hi, d.vertices[:0])
			texture.Begin()
			d.drawVertices(d.mesh.DrawMode, d.vertices)
			texture.End()
			lo, texture = hi, hitexture
		}
	}

	hi := batch.Len()
	d.vertices = instanceVertices(d.mesh, batch, lo, hi, d.vertices[:0])
	texture.Begin()
	d.drawVertices(d.mesh.DrawMode, d.vertices)
	texture.End()

	d.vslice.End()
}

func (d *Drawer) drawVertices(mode gl.Enum, verts []Vertex) {
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
