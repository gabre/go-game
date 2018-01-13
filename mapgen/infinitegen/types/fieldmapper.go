package types

import (
	"go-game/mapgen/infinitegen/noisegen"
	"engo.io/engo/common"
)

type Terrain = [][]float64

type FieldMapper interface {
	MapFields(noiseGen noisegen.NoiseGen, terrain Terrain, resolution int, tileset map[int]*common.Tile) (Layers []*common.TileLayer)
}
