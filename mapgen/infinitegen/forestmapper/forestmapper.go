package forestmapper

import (
	"go-game/mapgen/infinitegen/noisegen"
	"go-game/mapgen/infinitegen/types"
	"engo.io/engo/common"
	"engo.io/engo"
	"fmt"
)

const grassGid = 299
const treeGid = 1025

type forestMapper struct{}

func New() forestMapper {
	return forestMapper{}
}

func (f forestMapper) MapFields(noiseGen noisegen.NoiseGen, terrain types.Terrain, resolution int, startX int64, startZ int64, tileset map[int]*common.Tile) (Layers types.LayeredTiles, Width float32, Height float32) {
	fmt.Printf("MAPFIELDS: %d %d", startX, startZ)
	grass := make([][]*common.Tile, len(terrain))
	bushes := make([][]*common.Tile, len(terrain))
	trees := make([][]*common.Tile, len(terrain))

	uniformW := tileset[grassGid].Image.Width()
	uniformH := tileset[grassGid].Image.Height()
	startXf := float32(startX * int64(resolution)) * uniformW
	startZf := float32(startZ * int64(resolution)) * uniformH

	for col, _ := range terrain {
		grass[col] = make([]*common.Tile, len(terrain[col]))
		bushes[col] = make([]*common.Tile, len(terrain[col]))
		trees[col] = make([]*common.Tile, len(terrain[col]))
		for row, _ := range terrain {
			coord := engo.Point{startXf + float32(col) * uniformW, startZf + float32(row) * uniformH}
			grass[col][row] = mapGrass(coord, terrain[col][row], tileset)
			bushes[col][row] = mapBush(coord, terrain[col][row])
			trees[col][row] = mapTree(coord, terrain[col][row], tileset)
		}
	}

	return types.LayeredTiles{grass, bushes, trees}, uniformW * float32(resolution), uniformH * float32(resolution)
}

func mapGrass(coord engo.Point, value float64, tileset map[int]*common.Tile) *common.Tile {
	return &common.Tile{coord, tileset[grassGid].Image}
}

func mapBush(coord engo.Point, value float64) *common.Tile {
	return nil
}

func mapTree(coord engo.Point, value float64, tileset map[int]*common.Tile) *common.Tile {
	if value > 0.06 {
		return &common.Tile{coord, tileset[treeGid].Image}
	}
	return nil
}
