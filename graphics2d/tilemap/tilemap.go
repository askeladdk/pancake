package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
)

// TileID identifies a specific tile in a TileSet.
type TileID int

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

	// Metrics returns the metrics.
	Metrics() Metrics
}
