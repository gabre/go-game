package engine

import (
	"go-game/mapgen"
	"engo.io/ecs"
	"engo.io/engo/common"
	"log"
	"image/color"
	m "go-game/engine/map"
)
const definedR = int64(1) // [ -r ... 1 .. r ] tiles are rendered in each direction (horizontally/vertically)
                          // for renderNearestMapParts

type GameEngine struct {
	resolution      int
	mapGenFactory   func(int) (mapgen.MapGenerator, error)
	mapGenerator    mapgen.MapGenerator
	camZoom         float64
	renderSystem    *common.RenderSystem
	camSystem       *common.CameraSystem
}

func NewGameEngine(resolution int, generatorFactory func(int)(mapgen.MapGenerator, error)) *GameEngine {
	return &GameEngine{resolution: resolution, mapGenFactory: generatorFactory, camZoom: 1.0}
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
	chSys := m.NewChunkSystem(e.camSystem, e.renderSystem, e.mapGenerator, definedR)
	w.AddSystem(chSys)

	//engo.Mailbox.Listen("CameraMessage", func(msg engo.Message) {
	//	 _, ok := msg.(common.CameraMessage)
	//	 if !ok {
	//	    return
	//	 }
	//})
	chSys.Init()
}