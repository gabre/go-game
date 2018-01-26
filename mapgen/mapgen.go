package mapgen

import (
	"engo.io/engo/common"
	"engo.io/ecs"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Level struct {
	Tiles  [][][]*Tile
	Width  float32
	Height float32
}

type MapGenerator interface {
	GenerateMap(x int64, z int64) (Level, error)
}
