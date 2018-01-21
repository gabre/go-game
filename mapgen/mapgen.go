package mapgen

import (
	"engo.io/engo/common"
)

type Level struct {
	Tiles  [][][]*common.Tile
	Width  float32
	Height float32
}

type MapGenerator interface {
	GenerateMap(x int64, z int64) (Level, error)
}
