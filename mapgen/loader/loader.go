package loader

import (
	"log"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/tmx"
	// if this is NOT imported, the image format won't be recognized
	_ "image/png"
)

const defaultMapFile = "data/test_base64.tmx"

type MapLoader struct{}

func (m MapLoader) GenerateMap() (*tmx.Map, map[string]map[string]*gfx.Object) {
	tmxMap, layers, err := tmx.LoadFile(defaultMapFile, nil)
	if err != nil {
		log.Panicf("Error when loading map: %s", err)
	}
	return tmxMap, layers
}
