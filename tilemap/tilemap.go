package tilemap

import (
	"image"
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// TileID identifies a specific tile in a TileSet.
type TileID uint32

// Absent denotes that a tile does not exist.
const Absent TileID = 0xFFFFFFFF

// TileSet is a list of images that share a common underlying texture.
type TileSet interface {
	// Texture returns the underlying texture.
	Texture() *graphics.Texture

	// TileRegion returns the TextureRegion of the given tile.
	TileRegion(tileId TileID) graphics.TextureRegion

	// SameAutoTile returns whether two tiles belong to the same AutoTiler.
	SameAutoTile(tileId0, tileId1 TileID) bool

	// IsAutoTile returns whether the tile supports autotiling.
	IsAutoTile(tileId TileID) (base TileID, autoTiler AutoTiler, ok bool)
}

// TileMap is a two-dimensional grid of tiles.
type TileMap interface {
	// TileAt returns the TileID at a grid position.
	TileAt(cell Coordinate) TileID

	// SetTileAt sets the TileID at a grid position.
	SetTileAt(cell Coordinate, tileId TileID)

	// TintColorAt reports the tint color of a tile at a grid position.
	TintColorAt(cell Coordinate) color.NRGBA

	// TileSet returns the TileSet.
	TileSet() TileSet

	// Bounds returns the bounds of the TileMap measured in pixels.
	Bounds() image.Rectangle

	// RangeTilesInViewport iterates over all tiles in the viewport
	// and returns their positions and modelviews.
	RangeTilesInViewport(viewport image.Rectangle, fn func(cell Coordinate, modelview mathx.Aff3))
}
