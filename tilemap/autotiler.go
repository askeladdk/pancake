package tilemap

import "image"

// AutoTiler maps a tile neighbour bitset to a TileID.
type AutoTiler interface {
	AutoTile(bitset uint8) TileID
}

var neighbours = []image.Point{
	{+1, -1}, // NE
	{+1, +1}, // SE
	{-1, +1}, // SW
	{-1, -1}, // NW
	{+0, -1}, // N
	{+1, +0}, // E
	{+0, +1}, // S
	{-1, +0}, // W
}

// AutoTile modifies the TileMap by autotiling all tiles inside the given area.
func AutoTile(tileMap TileMap, cellBounds image.Rectangle) {
	tileSet := tileMap.TileSet()
	for y := cellBounds.Min.Y; y < cellBounds.Max.Y; y++ {
		for x := cellBounds.Min.X; x < cellBounds.Max.X; x++ {
			if id0 := tileMap.TileAt(x, y); id0 != Absent {
				if base, autoTiler := tileSet.IsAutoTile(id0); autoTiler != nil {
					var bitset uint8
					for bit, neighbour := range neighbours {
						id1 := tileMap.TileAt(x+neighbour.X, y+neighbour.Y)
						if id1 == Absent || !tileSet.SameBaseTile(id0, id1) {
							bitset |= 1 << bit
						}
					}
					tileMap.SetTileAt(x, y, base+autoTiler.AutoTile(bitset))
				}
			}
		}
	}
}

type arrayAutoTiler [256]TileID

func (a *arrayAutoTiler) AutoTile(bitset uint8) TileID { return a[bitset] }

// BlobTiler is an AutoTiler that maps corners and edges to 47 unique tiles.
var BlobTiler = AutoTiler(&arrayAutoTiler{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
	0x10, 0x10, 0x11, 0x11, 0x12, 0x12, 0x13, 0x13, 0x10, 0x10, 0x11, 0x11, 0x12, 0x12, 0x13, 0x13,
	0x14, 0x14, 0x14, 0x14, 0x15, 0x15, 0x15, 0x15, 0x16, 0x16, 0x16, 0x16, 0x17, 0x17, 0x17, 0x17,
	0x18, 0x18, 0x18, 0x18, 0x19, 0x19, 0x19, 0x19, 0x18, 0x18, 0x18, 0x18, 0x19, 0x19, 0x19, 0x19,
	0x1A, 0x1B, 0x1A, 0x1B, 0x1A, 0x1B, 0x1A, 0x1B, 0x1C, 0x1D, 0x1C, 0x1D, 0x1C, 0x1D, 0x1C, 0x1D,
	0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E,
	0x1F, 0x1F, 0x1F, 0x1F, 0x1F, 0x1F, 0x1F, 0x1F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
	0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21, 0x21,
	0x22, 0x23, 0x24, 0x25, 0x22, 0x23, 0x24, 0x25, 0x22, 0x23, 0x24, 0x25, 0x22, 0x23, 0x24, 0x25,
	0x26, 0x26, 0x27, 0x27, 0x26, 0x26, 0x27, 0x27, 0x26, 0x26, 0x27, 0x27, 0x26, 0x26, 0x27, 0x27,
	0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28, 0x28,
	0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29, 0x29,
	0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B, 0x2A, 0x2B,
	0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C,
	0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D,
	0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E, 0x2E,
})

// EdgeTiler is an autotiler that maps edges to 16 unique tiles.
var EdgeTiler = AutoTiler(&arrayAutoTiler{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
	0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
	0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03,
	0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
	0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
	0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09, 0x09,
	0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A, 0x0A,
	0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B,
	0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C,
	0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D, 0x0D,
	0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E, 0x0E,
	0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
})
