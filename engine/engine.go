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
const camZoomSpeed = 0.01 // 0.01x zoom for each scroll wheel click.
const camMinZoom = 0.1

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type GameEngine struct {
	resolution      int
	mapGenFactory   func(int) (mapgen.MapGenerator, error)
	mapGenerator    mapgen.MapGenerator
	camZoom         float64
	renderSystem    *common.RenderSystem
	camSystem       *common.CameraSystem

	// TODO The following four are surely redundant.
	// e.g. mapTopLeftCoord could be used with lastMapTopLeftCoord
	// The following are measured in integer X Y coordinates (which chunk)
	// Currently rendered chunkset's CENTER chunk's top left point's coordinates
	viewPoint       point
	// Last rendered chunkset's CENTER chunk's top left point's coordinates (equals `viewPoint` if everything is already rendered)
	lastViewPoint   point
	// The following are measured in "camera pixels":
	// Camera (which can be "scrolled away" from `viewPoint`) top left coordinates
	camTopLeftCoord engo.Point
	// Actual chunkset's (map) top left point's coordinates
	mapTopLeftCoord engo.Point

	cache           map[engo.Point]bool

	chunksetWidth  float32
	chunksetHeight float32

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
	c := make(map[engo.Point]bool)
	return &GameEngine{resolution: resolution, mapGenFactory: generatorFactory, camZoom: 1.0, viewPoint: p, lastViewPoint: p, cache: c}
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
			e.camTopLeftCoord = engo.Point{e.camSystem.X(), e.camSystem.Y()}
			e.mapTopLeftCoord = engo.Point{0, 0}
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
}

func (e *GameEngine) renderNearestMapPartsIfNeeded() {
	if (e.viewPoint != e.lastViewPoint) {
		e.renderNearestMapParts()
		e.lastViewPoint = e.viewPoint
	}
}

func (e *GameEngine) renderNearestMapParts() {
	fmt.Printf("-----  %d %d\n", e.viewPoint.x, e.viewPoint.y)
	// TODO this is not OK, adjust length
	chunks := make([]levelWithCoords, 0)
	for x := (e.viewPoint.x - r); x <= (e.viewPoint.x + r); x++ {
		for y := (e.viewPoint.y - r); y <= (e.viewPoint.y + r); y++ {
			p := engo.Point{float32(x), float32(y)}
			_, ok := e.cache[p]
			if (!ok) {
				fmt.Printf("RENDERING: %d %d  (for %d %d)\n", x, y, e.viewPoint.x, e.viewPoint.y)
				e.cache[p] = true
				chunks = append(chunks, levelWithCoords{p, e.generateMapCoords(x, y)})
			}
		}
	}

	// TODO we zero chunksetWidth chunksetHeight here as we deleted the previous map (which is not true)
	e.chunksetWidth = 0.0
	e.chunksetHeight = 0.0
	e.render(chunks...)
}

func (e *GameEngine) generateMapCoords(x int64, z int64) mapgen.Level {
	levelData, err := e.mapGenerator.GenerateMap(x, z)
	if err != nil {
		log.Panicf("Error while generating map: %s", err)
	}
	return levelData
}

func (e *GameEngine) render(chunks ...levelWithCoords) {
	if (len(chunks) == 0) {
		return
	}
	// TODO these should be in some kind of init
	e.chunkHeight = chunks[0].level.Height
	e.chunkWidth = chunks[0].level.Width
	numOfChunksPerRow := float32(r + 1 + r)
	e.chunksetHeight = chunks[0].level.Height * numOfChunksPerRow
	e.chunksetWidth = chunks[0].level.Width * numOfChunksPerRow
	fmt.Printf("W H %d %d\n", e.chunkWidth, e.chunkHeight)

	// Create render and space components for each of the tiles in all layers
	tileComponents := make([]*Tile, 0)
	for layer := range chunks[0].level.Tiles {
		for col := range chunks[0].level.Tiles[0] {
			for _, chunk := range chunks {
				// Keeping track of full map size

				for row := range chunks[0].level.Tiles[0][0] {
					tileElement := chunk.level.Tiles[layer][row][col]
					if tileElement != nil && tileElement.Image != nil {
						tile := &Tile{BasicEntity: ecs.NewBasic()}
						tile.RenderComponent = common.RenderComponent{
							Drawable: tileElement,
							Scale:    engo.Point{1, 1},
						}
						tile.SpaceComponent = common.SpaceComponent{
							Position: tileElement.Point,
							Width:    0,
							Height:   0,
						}
						tileComponents = append(tileComponents, tile)
					}
				}
			}
		}
	}

	// Add each of the tiles entities and its components to the render system
	for _, v := range tileComponents {
		e.renderSystem.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
	}
}