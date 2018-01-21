package types

import (
	"go-game/mapgen/infinitegen/noisegen"
	"engo.io/engo/common"
)

type Terrain = [][]float64
type LayeredTiles = [][][]*common.Tile

type FieldMapper interface {
	MapFields(noiseGen noisegen.NoiseGen, terrain Terrain, resolution int, startX int64, startZ int64, tileset map[int]*common.Tile) (Layers LayeredTiles, Width float32, Height float32)
}
