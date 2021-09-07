package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/mathx"
)

// TileID identifies a specific tile in a TileSet.
type TileID uint32

// Absent denotes that a tile does not exist.
const Absent TileID = 0xFFFFFFFF

// TileSet is a set of images that share a common underlying texture.
type TileSet interface {
	// Texture returns the underlying texture.
	Texture() *pancake.Texture

	// TileRegion returns the TextureRegion of the given tile.
	TileRegion(TileID) pancake.TextureRegion

	// SameBaseTile reports whether two tiles share the same base tile.
	SameBaseTile(id0, id1 TileID) bool

	// IsAutoTile returns whether the tile supports autotiling.
	IsAutoTile(TileID) (base TileID, autoTiler AutoTiler)

	// TileSize reports the size of the tiles in pixels.
	TileSize() mathx.Vec2
}

// TileMap is a two-dimensional grid of tiles.
type TileMap interface {
	// TileAt returns the TileID at a grid position.
	TileAt(x, y int) TileID

	// SetTileAt sets the TileID at a grid position.
	SetTileAt(x, y int, id TileID)

	// TintColorAt reports the tint color of a tile at a grid position.
	TintColorAt(x, y int) color.Color

	// TileSet returns the TileSet.
	TileSet() TileSet
}
