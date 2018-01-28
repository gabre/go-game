package _map

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"engo.io/engo"
	"engo.io/engo/math"
	"go-game/mapgen"
	"log"
)

// TODO move into fields
var zeroX, zeroY = float32(0), float32(0)

type levelWithCoords struct {
	coords engo.Point
	level  mapgen.Level
}

type point struct {
	x, y int64
}

type ChunkSystem struct {
	cameraSys *common.CameraSystem
	renderSys *common.RenderSystem
	mapGenerator mapgen.MapGenerator
	r int64 // number of chunks rendered vertically/horizontally counted from the central chunk
			// e.g. [ -r ... CENTRAL ... +r ]
	init            bool

	// Currently rendered chunkset's CENTER chunk's top left point's coordinates
	viewPoint       point
	// Last rendered chunkset's CENTER chunk's top left point's coordinates (equals `viewPoint` if everything is already rendered)
	lastViewPoint   point

	cache           map[engo.Point]mapgen.Level

	chunkWidth     float32
	chunkHeight    float32
}

func NewChunkSystem(cameraSys *common.CameraSystem, renderSys *common.RenderSystem, mapGenerator mapgen.MapGenerator, r int64) *ChunkSystem {
	p := point{0,0}
	c := make(map[engo.Point]mapgen.Level)
	zeroX = cameraSys.X()
	zeroY = cameraSys.Y()
	return &ChunkSystem{cameraSys: cameraSys, renderSys: renderSys, mapGenerator:mapGenerator, r: r, viewPoint: p, lastViewPoint: p, cache: c, init: true}
}

func (r *ChunkSystem) Init() {
	if (r.init) {
		r.renderNearestMapParts()
		r.init = false
	}
}

// Remove is called whenever an Entity is removed from the World, in order to remove it from this sytem as well
func (*ChunkSystem) Remove(ecs.BasicEntity) {}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
// TODO this should be a method of engine
func (r *ChunkSystem) Update(dt float32) {
	r.UpdateViewPointWithCamera()
	r.renderNearestMapPartsIfNeeded()
}
func (r *ChunkSystem) UpdateViewPointWithCamera() {
	halfWW := engo.WindowWidth() / 2
	halfWH := engo.WindowHeight() / 2
	// zero is (zeroX - halfW), so camera X measured from this is camX
	camX := r.cameraSys.X() + (zeroX - halfWW)
	camY := r.cameraSys.Y() + (zeroY - halfWH)
	actualViewPointX := camScalarToDiscretePoint(camX, r.chunkWidth)
	actualViewPointY := camScalarToDiscretePoint(camY, r.chunkHeight)
	r.viewPoint = point{actualViewPointX, actualViewPointY}
}

func (r *ChunkSystem) numOfChunksPerRow() int64 {
	return (r.r + 1 + r.r)
}

func camScalarToDiscretePoint(camScalar float32, unitLength float32) int64 {
	translation := unitLength / 2
	camScalar += translation
	units := camScalar / unitLength
	return int64(math.Floor(units))
}

func (e *ChunkSystem) renderNearestMapPartsIfNeeded() {
	if (e.viewPoint != e.lastViewPoint) {
		e.renderNearestMapParts()
		e.lastViewPoint = e.viewPoint
	}
}

func (e *ChunkSystem) renderNearestMapParts() {
	toRemove := make([]levelWithCoords, 0, int(e.numOfChunksPerRow() * e.numOfChunksPerRow()))
	chunks := make([]levelWithCoords, 0, int(e.numOfChunksPerRow() * e.numOfChunksPerRow()))
	for i := -e.r; i <= e.r; i++ {
		for j := -e.r; j <= e.r; j++ {
			x := e.viewPoint.x + j
			y := e.viewPoint.y + i
			lx := e.lastViewPoint.x + j
			ly := e.lastViewPoint.y + i
			p := engo.Point{float32(x), float32(y)}
			lp := engo.Point{float32(lx), float32(ly)}
			lvl, ok := e.cache[p]
			if (!ok) {
				lvl = e.generateMapCoords(x, y)
				e.cache[p] = lvl
			}
			oldlvl, oldWasCached := e.cache[lp]
			if (!e.init && oldWasCached) {
				toRemove = append(toRemove, levelWithCoords{lp, oldlvl})
			}
			chunks = append(chunks, levelWithCoords{p, lvl})
		}
	}

	e.renderAndDeletePrevious(chunks, toRemove)
}

func (e *ChunkSystem) generateMapCoords(x int64, z int64) mapgen.Level {
	levelData, err := e.mapGenerator.GenerateMap(x, z)
	if err != nil {
		log.Panicf("Error while generating map: %s", err)
	}
	return levelData
}

// This function renders a bunch of chunks and deletes the previous (which has equal size)
// len(chunks) must be equal to len(toRemove) or toRemove must be nil
func (e *ChunkSystem) renderAndDeletePrevious(chunks []levelWithCoords, toRemove []levelWithCoords) {
	if (len(chunks) == 0) {
		return
	}
	if (e.init) {
		e.chunkHeight = chunks[0].level.Height
		e.chunkWidth = chunks[0].level.Width
	}

	// Create render and space components for each of the tiles in all layers
	tiles := make([]*mapgen.Tile, 0, len(chunks) * len(chunks[0].level.Tiles) * len(chunks[0].level.Tiles[0]) * len(chunks[0].level.Tiles[0][0]))
	for layer := range chunks[0].level.Tiles {
		for chunkrow := int64(0); chunkrow < e.numOfChunksPerRow(); chunkrow++ {
			for row := range chunks[0].level.Tiles[0] {
				chunkrowStart := chunkrow * e.numOfChunksPerRow()
				chunkrowEnd := chunkrowStart + e.numOfChunksPerRow()
				for i := chunkrowStart; i < chunkrowEnd; i++ {
					for col := range chunks[0].level.Tiles[0][0] {
						if len(toRemove) != 0 {
							r := toRemove[i].level.Tiles[layer][row][col]
							if r != nil {
								e.renderSys.Remove(r.BasicEntity)
							}
						}
						v := chunks[i].level.Tiles[layer][row][col]
						// Add each of the tiles entities and its components to the render system
						if v != nil {
							tiles = append(tiles, v)
						}
					}
				}
			}
		}
	}
	for _, v := range tiles {
		e.renderSys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
	}
}