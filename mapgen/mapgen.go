package mapgen

import (
	"engo.io/engo/common"
)

type LayerNames = []string

type MapGenerator interface {
	GenerateMap(x int64, z int64) (*common.Level, error)
}
