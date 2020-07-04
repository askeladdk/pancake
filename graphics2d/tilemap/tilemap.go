package tilemap

import (
	"image"
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// TileID identifies a specific tile in a TileSet.
type TileID uint32

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
	TileAt(x, y int) TileID

	// SetTileAt sets the TileID at a grid position.
	SetTileAt(x, y int, tileId TileID)

	// TileColorAt gets the color of a tile at a grid position.
	TileColorAt(x, y int) color.NRGBA

	// TileSet returns the TileSet.
	TileSet() TileSet

	// Bounds returns the bounds of the TileMap measured in pixels.
	Bounds() image.Rectangle

	// AutoTileBitSet computes the bitset of a cell based on its eight neighbours.
	// The test function tests whether a neighbouring cell is of the same type as the centre cell.
	AutoTileBitSet(cx, cy int, testFunc func(x, y int) bool) uint8

	// RangeTilesInViewport iterates over all tiles in the viewport
	// and returns their positions and modelviews.
	RangeTilesInViewport(viewport image.Rectangle, fn func(x, y int, modelview mathx.Aff3))
}
