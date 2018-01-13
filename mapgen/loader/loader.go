package loader

import (
	// if this is NOT imported, the image format won't be recognized
	_ "image/png"
	"engo.io/engo/common"
	"go-game/util"
)

const defaultMapFile = "data/test_base64.tmx"

type MapLoader struct{}

func (m MapLoader) GenerateMap(x int64, z int64) (*common.Level, error) {
	return util.LoadTmxMap(defaultMapFile)
}
