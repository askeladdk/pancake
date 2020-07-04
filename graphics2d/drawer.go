package graphics2d

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/mathx"
)

var VertexFormat = graphics.AttribFormat{
	graphics.Vec2,  // XY
	graphics.Vec2,  // UV
	graphics.Byte4, // RGBA
}

type Vertex struct {
	XY   mathx.Vec2
	UV   mathx.Vec2
	RGBA color.NRGBA
}

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	DrawMode gl.Enum
}

type Batch interface {
	Len() int
	TintColorAt(i int) color.NRGBA
	Texture() *graphics.Texture
	TextureRegionAt(i int) graphics.TextureRegion
	ModelViewAt(i int) mathx.Aff3
	PivotAt(i int) mathx.Vec2
}

func MakeVertices(mesh Mesh, batch Batch, vertices []Vertex) []Vertex {
	tmpv := make([]Vertex, len(mesh.Vertices))

	for i := 0; i < batch.Len(); i++ {
		modelview := batch.ModelViewAt(i)
		region := batch.TextureRegionAt(i).Aff3()
		rgba := batch.TintColorAt(i)
		pivot := batch.PivotAt(i)
		for m, v := range mesh.Vertices {
			tmpv[m] = Vertex{
				XY: modelview.Project(v.XY.Add(pivot)),
				UV: region.Project(v.UV),
				RGBA: color.NRGBA{
					R: uint8(uint16(v.RGBA.R) * uint16(rgba.R) / 255),
					G: uint8(uint16(v.RGBA.G) * uint16(rgba.G) / 255),
					B: uint8(uint16(v.RGBA.B) * uint16(rgba.B) / 255),
					A: uint8(uint16(v.RGBA.A) * uint16(rgba.A) / 255),
				},
			}
		}
		for _, n := range mesh.Indices {
			vertices = append(vertices, tmpv[n])
		}
	}

	return vertices
}

type Drawer struct {
	verts   []Vertex
	vbuffer *graphics.Buffer
	vslice  *graphics.VertexSlice
	mesh    Mesh
}

func NewDrawer(maxinstances int, mesh Mesh) *Drawer {
	vbuffer := graphics.NewBuffer(VertexFormat, maxinstances*len(mesh.Indices), nil)
	vslice := graphics.NewVertexSlice(vbuffer)
	return &Drawer{
		vbuffer: vbuffer,
		vslice:  vslice,
		mesh:    mesh,
	}
}

func (d *Drawer) Draw(batches ...Batch) {
	if len(batches) > 1 {
		d.drawBatches(batches)
	} else if len(batches) == 1 {
		d.drawBatch(batches[0])
	}
}

func (d *Drawer) drawBatches(batches []Batch) {
	d.vslice.Begin()

	for i := 0; i < len(batches); {
		d.verts = MakeVertices(d.mesh, batches[i], d.verts[:0])
		texture := batches[i].Texture()

		for i = i + 1; i < len(batches); i++ {
			if batches[i].Texture() != texture {
				break
			}

			d.verts = MakeVertices(d.mesh, batches[i], d.verts)
		}

		texture.Begin()
		d.drawVertices(d.mesh.DrawMode, d.verts)
		texture.End()
	}

	d.vslice.End()
}

func (d *Drawer) drawBatch(batch Batch) {
	d.verts = MakeVertices(d.mesh, batch, d.verts[:0])
	texture := batch.Texture()

	d.vslice.Begin()
	texture.Begin()
	d.drawVertices(d.mesh.DrawMode, d.verts)
	texture.End()
	d.vslice.End()
}

func (d *Drawer) drawVertices(mode gl.Enum, verts []Vertex) {
	var indices []int
	for i := 0; i < len(verts); i += d.vslice.Len() {
		indices = append(indices, i)
	}
	indices = append(indices, len(verts))

	for i := 1; i < len(indices); i++ {
		lo := indices[i-1]
		hi := indices[i-0]
		vslice := d.vslice.Slice(0, hi-lo)
		vslice.SetData(verts[lo:hi])
		vslice.Draw(mode)
	}
}
