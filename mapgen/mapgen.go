package mapgen

import (
	"azul3d.org/engine/gfx"
	"azul3d.org/engine/tmx"
)

type MapGenerator interface {
	GenerateMap() (*tmx.Map, map[string]map[string]*gfx.Object)
}
