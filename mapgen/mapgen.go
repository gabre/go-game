package mapgen

import (
	"azul3d.org/engine/gfx"
)

type LayerNames = []string
type ObjectMap = map[string]map[string]*gfx.Object

type MapGenerator interface {
	GenerateMap(x int64, z int64) (ObjectMap, error)
}
