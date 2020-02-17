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
	Color(i int) color.NRGBA
	Texture() *graphics.Texture
	TextureRegion(i int) mathx.Aff3
	ModelView(i int) mathx.Aff3
}

func MakeVertices(mesh Mesh, batch Batch, vertices []Vertex) []Vertex {
	tmpv := make([]Vertex, len(mesh.Vertices))

	for i := 0; i < batch.Len(); i++ {
		modelview := batch.ModelView(i)
		region := batch.TextureRegion(i)
		rgba := batch.Color(i)
		for m := 0; m < len(tmpv); m++ {
			v := mesh.Vertices[m]
			tmpv[m] = Vertex{
				XY: modelview.Project(v.XY),
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

func (d *Drawer) DrawBatches(batches []Batch) {
	var verts []Vertex

	d.vslice.Begin()

	for i := 0; i < len(batches); {
		verts = MakeVertices(d.mesh, batches[i], verts[:0])
		texture := batches[i].Texture()

		for i = i + 1; i < len(batches); i++ {
			if batches[i].Texture() != texture {
				break
			}

			verts = MakeVertices(d.mesh, batches[i], verts)
		}

		texture.Begin()
		d.drawVertices(d.mesh.DrawMode, verts)
		texture.End()
	}

	d.vslice.End()
}

func (d *Drawer) DrawBatch(batch Batch) {
	verts := MakeVertices(d.mesh, batch, nil)
	texture := batch.Texture()

	d.vslice.Begin()
	texture.Begin()
	d.drawVertices(d.mesh.DrawMode, verts)
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
