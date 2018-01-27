package engine

import (
	"go-game/mapgen"
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"log"
	"image/color"
	"fmt"
)
const r = int64(1) // [ -r ... 1 .. r ] tiles are rendered in each direction (horizontally/vertically)
				   // for renderNearestMapParts
const numOfChunksPerRow = float32(r + 1 + r)

type GameEngine struct {
	resolution      int
	mapGenFactory   func(int) (mapgen.MapGenerator, error)
	mapGenerator    mapgen.MapGenerator
	camZoom         float64
	renderSystem    *common.RenderSystem
	camSystem       *common.CameraSystem

	init            bool

	// Currently rendered chunkset's CENTER chunk's top left point's coordinates
	viewPoint       point
	// Last rendered chunkset's CENTER chunk's top left point's coordinates (equals `viewPoint` if everything is already rendered)
	lastViewPoint   point

	cache           map[engo.Point]mapgen.Level

	chunkWidth     float32
	chunkHeight    float32
}

type levelWithCoords struct {
	coords engo.Point
	level  mapgen.Level
}

type point struct {
	x, y int64
}

func NewGameEngine(resolution int, generatorFactory func(int)(mapgen.MapGenerator, error)) *GameEngine {
	p := point{0,0}
	c := make(map[engo.Point]mapgen.Level)
	return &GameEngine{resolution: resolution, mapGenFactory: generatorFactory, camZoom: 1.0, viewPoint: p, lastViewPoint: p, cache: c, init: true}
}

func (e *GameEngine) Type() string {
	return "GameEngine"
}

func (e *GameEngine) Preload() {
	if (e.mapGenerator == nil) {
		var err error
		e.mapGenerator, err = e.mapGenFactory(e.resolution)
		if (err != nil) {
			log.Panicf("Map generation error: %s", err)
		}
	}
}

func (e *GameEngine) Setup(w *ecs.World) {
	common.SetBackground(color.Black)
	e.renderSystem = &common.RenderSystem{}
	w.AddSystem(e.renderSystem)
	w.AddSystem(&common.EdgeScroller{400, 20})
	for _, syst := range w.Systems() {
		switch sys := syst.(type) {
		case *common.CameraSystem:
			e.camSystem = sys
			println(e.camSystem.X(), e.camSystem.Y())
		}
	}
	chSys := NewChunkSystem(e, e.camSystem, r)
	w.AddSystem(chSys)

	//engo.Mailbox.Listen("CameraMessage", func(msg engo.Message) {
	//	 _, ok := msg.(common.CameraMessage)
	//	 if !ok {
	//	    return
	//	 }
	//})
	zeroX = e.camSystem.X()
	zeroY = e.camSystem.Y()
	e.renderNearestMapParts()
	e.init = false
}

func (e *GameEngine) renderNearestMapPartsIfNeeded() {
	if (e.viewPoint != e.lastViewPoint) {
		e.renderNearestMapParts()
		e.lastViewPoint = e.viewPoint
	}
}

func (e *GameEngine) renderNearestMapParts() {
	toRemove := make([]levelWithCoords, 0, numOfChunksPerRow * numOfChunksPerRow)
	chunks := make([]levelWithCoords, 0, numOfChunksPerRow * numOfChunksPerRow)
	for i := -r; i <= r; i++ {
		for j := -r; j <= r; j++ {
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

func (e *GameEngine) generateMapCoords(x int64, z int64) mapgen.Level {
	levelData, err := e.mapGenerator.GenerateMap(x, z)
	if err != nil {
		log.Panicf("Error while generating map: %s", err)
	}
	return levelData
}

// This function renders a bunch of chunks and deletes the previous (which has equal size)
// len(chunks) must be equal to len(toRemove) or toRemove must be nil
func (e *GameEngine) renderAndDeletePrevious(chunks []levelWithCoords, toRemove []levelWithCoords) {
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
		for chunkrow := 0; chunkrow < int(numOfChunksPerRow); chunkrow++ {
			for row := range chunks[0].level.Tiles[0] {
				chunkrowStart := chunkrow * int(numOfChunksPerRow)
				chunkrowEnd := chunkrowStart + int(numOfChunksPerRow)
				for i := chunkrowStart; i < chunkrowEnd; i++ {
					for col := range chunks[0].level.Tiles[0][0] {
						if len(toRemove) != 0 {
							r := toRemove[i].level.Tiles[layer][row][col]
							if r != nil {
								e.renderSystem.Remove(r.BasicEntity)
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
		e.renderSystem.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
	}
}