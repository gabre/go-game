package forestmapper

import (
	"go-game/mapgen/infinitegen/noisegen"
	"go-game/mapgen/infinitegen/types"
	"engo.io/engo/common"
	"engo.io/engo"
)

const grassGid = 299
const treeGid = 1025

type forestMapper struct{}

func New() forestMapper {
	return forestMapper{}
}

func (f forestMapper) MapFields(noiseGen noisegen.NoiseGen, terrain types.Terrain, resolution int, tileset map[int]*common.Tile) (Layers []*common.TileLayer) {
	grass := make([]*common.Tile, 0)
	bushes := make([]*common.Tile, 0)
	trees := make([]*common.Tile, 0)
	for row, _ := range terrain {
		for col, _ := range terrain {
			uniformW := tileset[grassGid].Image.Width()
			uniformH := tileset[grassGid].Image.Height()
			coord := engo.Point{float32(col) * uniformW, float32(row) * uniformH}
			grass = mapGrass(coord, terrain[col][row], grass, tileset)
			bushes = mapBush(coord, terrain[col][row], bushes)
			trees = mapTree(coord, terrain[col][row], trees, tileset)
			// get the tilesheets in order and in generic format
			// sort.Sort(common.ByFirstgid(tmxLevel.Tilesets))
			// ts := make([]*tilesheet, len(tmxLevel.Tilesets))
			// for i, tts := range tmxLevel.Tilesets {
			//	ts[i] = &tilesheet{tts.Image, tts.Firstgid}
			// }
		}
	}
	layerGrass := common.TileLayer{Name: "grass", Tiles: grass}
	layerBush := common.TileLayer{Name: "bushes", Tiles: bushes}
	layerTree := common.TileLayer{Name: "trees", Tiles: trees}
	return []*common.TileLayer{&layerGrass, &layerTree, &layerBush}
}

func mapGrass(coord engo.Point, value float64, a []*common.Tile, tileset map[int]*common.Tile) []*common.Tile {
	return append(a, &common.Tile{coord, tileset[grassGid].Image})
}

func mapBush(coord engo.Point, value float64, a []*common.Tile) []*common.Tile {
	return a
}

func mapTree(coord engo.Point, value float64, a []*common.Tile, tileset map[int]*common.Tile) []*common.Tile {
	if value > 0.06 {
		a = append(a, &common.Tile{coord, tileset[treeGid].Image})
	}
	return a
}
