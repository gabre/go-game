package types

import (
	"go-game/mapgen/infinitegen/noisegen"
	"engo.io/engo/common"
	"go-game/mapgen"
)

type Terrain = [][]float64
type LayeredTiles = [][][]*mapgen.Tile

type FieldMapper interface {
	MapFields(noiseGen noisegen.NoiseGen, terrain Terrain, resolution int, startX int64, startZ int64, tileset map[int]*common.Tile) (Layers LayeredTiles, Width float32, Height float32)
}
