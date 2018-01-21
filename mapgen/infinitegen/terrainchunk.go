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
	// x and z are tile "ID numbers" - we need to convert them to corner point coordinates
	fx := float64(x * int64(resolution))
	fz := float64(z * int64(resolution))
	terrain := make(types.Terrain, res)
	for ix := int64(0); ix < res; ix++  {
		terrain[ix] = make([]float64, res)
		for iz := int64(0); iz < res; iz++ {
			noiseX := float64(fx) + float64(ix) / float64(res)
			noiseZ := float64(fz) + float64(iz) / float64(res)
			terrain[ix][iz] = noiseGen.GetNoise(noiseX, noiseZ)
		}
	}
	return &terrainChunk{x:x, z:z, noiseGen:noiseGen, resolution: resolution, terrain: terrain}
}
