package loader

import (
	"azul3d.org/engine/tmx"
	// if this is NOT imported, the image format won't be recognized
	_ "image/png"
	"go-game/mapgen"
)

const defaultMapFile = "data/test_base64.tmx"

type MapLoader struct{}

func (m MapLoader) GenerateMap(x int64, z int64) (mapgen.ObjectMap, error) {
	tmxMap, layers, err := tmx.LoadFile(defaultMapFile, nil)
	if err != nil {
		return nil, err
	}
	finalLayers := make(mapgen.ObjectMap)
	for _, layer := range tmxMap.Layers {
		l, ok := layers[layer.Name]
		if ok {
			finalLayers[layer.Name] = l
		}
	}
	return finalLayers, nil
}
