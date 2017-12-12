package infinitegen

import (
	"go-game/mapgen/infinitegen/noisegen"
	"go-game/mapgen/infinitegen/types"
)

type terrainChunk struct {
	x          int64
	z          int64
	noiseGen   noisegen.NoiseGen
	resolution int
	terrain    types.Terrain
}

func newChunk(x int64, z int64, noiseGen noisegen.NoiseGen, resolution int) *terrainChunk {
	res := int64(resolution)
	terrain := make(types.Terrain, res)
	for ix := int64(0); ix < res; ix++  {
		terrain[ix] = make([]float64, res)
		for iz := int64(0); iz < res; iz++ {
			noiseX := float64(x) + float64(ix) / float64(res)
			noiseZ := float64(z) + float64(iz) / float64(res)
			terrain[ix][iz] = noiseGen.GetNoise(noiseX, noiseZ)
		}
	}
	return &terrainChunk{x:x, z:z, noiseGen:noiseGen, resolution: resolution, terrain: terrain}
}
