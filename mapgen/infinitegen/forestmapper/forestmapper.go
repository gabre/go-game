package forestmapper

import (
	"go-game/mapgen/infinitegen/noisegen"
	"azul3d.org/engine/tmx"
	"go-game/mapgen/infinitegen/types"
	"log"
)

const grassGid = 299
const treeGid = 1025

type forestMapper struct {}

func New() forestMapper {
	return forestMapper{}
}

func (f forestMapper) MapFields(noiseGen noisegen.NoiseGen, terrain types.Terrain, resolution int) (Layers []*tmx.Layer) {
	grass := make(map[tmx.Coord]uint32)
	bushes := make(map[tmx.Coord]uint32)
	trees := make(map[tmx.Coord]uint32)
	for x, row := range terrain {
		for y, value := range row {
			coord := tmx.Coord{x, y}
			log.Print("....")
			mapGrass(coord, value, grass)
			mapBush(coord, value, bushes)
			mapTree(coord, value, trees)
		}
	}
	layerGrass := tmx.Layer{Name:"grass", Opacity:0, Visible:true, Tiles:grass}
	layerBush := tmx.Layer{Name:"bushes", Opacity:0, Visible:true, Tiles:bushes}
	layerTree := tmx.Layer{Name:"trees", Opacity:0, Visible:true, Tiles:trees}
	log.Print(layerTree)
	return []*tmx.Layer{&layerGrass, &layerTree, &layerBush}
}

func mapGrass(coord tmx.Coord, value float64, m map[tmx.Coord]uint32) {
	m[coord] = grassGid
}

func mapBush(coord tmx.Coord, value float64, m map[tmx.Coord]uint32) {

}

func mapTree(coord tmx.Coord, value float64, m map[tmx.Coord]uint32) {
	log.Printf("Value: %f,   %f", value, noisegen.GetRangeMax())
	if value > 0.1 {
		log.Print("Beep!")
		m[coord] = treeGid
	}
}