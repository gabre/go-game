package infinitegen

import (
	"go-game/mapgen/infinitegen/forestmapper"
	"go-game/mapgen/infinitegen/noisegen"
	"go-game/mapgen/infinitegen/types"

	"engo.io/engo/common"
	"fmt"
	"go-game/util"
	_ "image/png"
)

const emptyMapPath = "data/empty.tmx"
const tilesDir = "data/tiles"

type objCoord struct {
	x int64
	z int64
}
type levelMap = map[objCoord]*common.Level

type MapLoader struct {
	noiseGen   noisegen.NoiseGen
	resolution int
	chunks     levelMap
	mapper     types.FieldMapper
	initialMap *common.Level
}

func New(seed int64, resolution int) (*MapLoader, error) {
	chunks := make(levelMap)
	noiseGen := noisegen.New(seed)
	var fieldMapper types.FieldMapper = forestmapper.New()
	// 1) we create an empty *Map (=: m)
	initialMap, err := util.LoadTmxMap(emptyMapPath)
	if err != nil {
		return nil, fmt.Errorf("Loading of initial map failed: %s", err)
	}
	initialMap.TileHeight = resolution
	initialMap.TileWidth = resolution
	if err != nil {
		return nil, fmt.Errorf("Infinite generator initial map problem: %s", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Infinite generator image loading problem: %s", err)
	}
	return &MapLoader{chunks: chunks, noiseGen: noiseGen, resolution: resolution, mapper: fieldMapper, initialMap: initialMap}, nil
}

func (m *MapLoader) GenerateMap(x int64, z int64) (*common.Level, error) {
	coord := objCoord{x, z}
	objMap, ok := m.chunks[coord]
	if !ok {
		chunk := newChunk(x, z, m.noiseGen, m.resolution)
		// 2) we generate the layers using a mapper
		// 3) we add the layers to m
		// Layers are always changed
		newMap := *m.initialMap
		newMap.TileLayers = m.mapper.MapFields(m.noiseGen, chunk.terrain, chunk.resolution, newMap.Tileset)
		objMap = &newMap
		m.chunks[coord] = objMap
	}

	return objMap, nil
}