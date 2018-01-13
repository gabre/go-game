package noisegen

import (
	"github.com/aquilax/go-perlin"
	"math"
)

const iterations int = 1
const alpha float64 = 2
const beta float64 = 2

type NoiseGen struct {
	perlinGen *perlin.Perlin
}

func New(seed int64) NoiseGen {
	p := perlin.NewPerlin(alpha, beta, iterations, seed)
	return NoiseGen{p}
}

func (p NoiseGen) GetNoise(x float64, z float64) float64 {
	v := p.perlinGen.Noise2D(x, z)
	return v
}

func GetRangeMin() float64 {
	return -1 * GetRangeMax()
}

func GetRangeMax() float64 {
	return math.Sqrt2 / 2
}