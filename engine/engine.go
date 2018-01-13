package engine

import (
	"go-game/mapgen"
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"log"
	"image/color"
)

const camZoomSpeed = 0.01 // 0.01x zoom for each scroll wheel click.
const camMinZoom = 0.1

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type GameEngine struct {
	mapGenFactory func() (mapgen.MapGenerator, error)
	mapGenerator  mapgen.MapGenerator
	camZoom       float64
	renderSystem  *common.RenderSystem
}

func NewGameEngine(generatorFactory func()(mapgen.MapGenerator, error)) *GameEngine {
	return &GameEngine{mapGenFactory: generatorFactory, camZoom: 1.0}
}

func (e *GameEngine) Type() string {
	return "GameEngine"
}

func (e *GameEngine) Preload() {
	if (e.mapGenerator == nil) {
		var err error
		e.mapGenerator, err = e.mapGenFactory()
		if (err != nil) {
			log.Panicf("Map generation error: %s", err)
		}
	}
}

func (e *GameEngine) Setup(w *ecs.World) {
	common.SetBackground(color.White)
	w.AddSystem(&common.RenderSystem{})
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			e.renderSystem = sys
		}
	}

	e.renderMapCoords(0,0)
}

func (e *GameEngine) renderMapCoords(x int64, z int64) {
	levelData, err := e.mapGenerator.GenerateMap(x, z)
	if err != nil {
		log.Panicf("Error while generating map: %s", err)
	}
	e.renderTiles(levelData)
}

func (e *GameEngine) renderTiles(levelData *common.Level) {
	// Create render and space components for each of the tiles in all layers
	tileComponents := make([]*Tile, 0)

	for _, tileLayer := range levelData.TileLayers {

		for _, tileElement := range tileLayer.Tiles {

			if tileElement.Image != nil {

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

	// Do the same for all image layers
	for _, imageLayer := range levelData.ImageLayers {
		for _, imageElement := range imageLayer.Images {
			if imageElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: imageElement,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: imageElement.Point,
					Width:    0,
					Height:   0,
				}

				tileComponents = append(tileComponents, tile)
			}
		}
	}

	// Add each of the tiles entities and its components to the render system
	for _, v := range tileComponents {
		e.renderSystem.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
	}
}