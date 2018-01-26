package forestmapper

import (
	"go-game/mapgen/infinitegen/noisegen"
	"go-game/mapgen/infinitegen/types"
	"engo.io/engo/common"
	"engo.io/engo"
	"fmt"
	"go-game/mapgen"
	"engo.io/ecs"
)

const grassGid = 299
const treeGid = 1025

type forestMapper struct{}

func New() forestMapper {
	return forestMapper{}
}

func (f forestMapper) MapFields(noiseGen noisegen.NoiseGen, terrain types.Terrain, resolution int, startX int64, startZ int64, tileset map[int]*common.Tile) (Layers types.LayeredTiles, Width float32, Height float32) {
	fmt.Printf("MAPFIELDS: %d %d", startX, startZ)
	grass := make([][]*mapgen.Tile, len(terrain))
	bushes := make([][]*mapgen.Tile, len(terrain))
	trees := make([][]*mapgen.Tile, len(terrain))

	uniformW := tileset[grassGid].Image.Width()
	uniformH := tileset[grassGid].Image.Height()
	startXf := float32(startX * int64(resolution)) * uniformW
	startZf := float32(startZ * int64(resolution)) * uniformH

	for row, _ := range terrain {
		grass[row] = make([]*mapgen.Tile, len(terrain[row]))
		bushes[row] = make([]*mapgen.Tile, len(terrain[row]))
		trees[row] = make([]*mapgen.Tile, len(terrain[row]))
		for col, _ := range terrain {
			coord := engo.Point{startXf + float32(col) * uniformW, startZf + float32(row) * uniformH}
			grass[row][col] = mapGrass(coord, terrain[row][col], tileset)
			bushes[row][col] = mapBush(coord, terrain[row][col])
			trees[row][col] = mapTree(coord, terrain[row][col], tileset, resolution * resolution)
		}
	}

	return types.LayeredTiles{grass, bushes, trees}, uniformW * float32(resolution), uniformH * float32(resolution)
}

func mapGrass(coord engo.Point, value float64, tileset map[int]*common.Tile) *mapgen.Tile {
	return createTile(&common.Tile{coord, tileset[grassGid].Image}, 0)
}

func mapBush(coord engo.Point, value float64) *mapgen.Tile {
	return nil
}

func mapTree(coord engo.Point, value float64, tileset map[int]*common.Tile, extra int) *mapgen.Tile {
	if value > 0.1 {
		return createTile(&common.Tile{coord, tileset[treeGid].Image}, float32(extra + 10))
	}
	return nil
}

func createTile(tileElement *common.Tile, zIndex float32) *mapgen.Tile {
	tile := &mapgen.Tile{BasicEntity: ecs.NewBasic()}
	tile.RenderComponent = common.RenderComponent{
		Drawable: tileElement,
		Scale:    engo.Point{1, 1},
	}
	tile.SpaceComponent = common.SpaceComponent{
		Position: tileElement.Point,
		Width:    0,
		Height:   0,
	}

	// tile.RenderComponent.SetZIndex(tileElement.Y + zIndex)
	return tile
}