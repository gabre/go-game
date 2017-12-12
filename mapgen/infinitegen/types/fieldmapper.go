package types

import (
	"go-game/mapgen/infinitegen/noisegen"
	"azul3d.org/engine/tmx"
)

type Terrain = [][]float64

type FieldMapper interface {
	MapFields(noiseGen noisegen.NoiseGen, terrain Terrain, resolution int) (Layers []*tmx.Layer)
}