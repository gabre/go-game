package main

import (
	"go-game/engine"

	"azul3d.org/engine/gfx/window"
	"go-game/mapgen/infinitegen"
	"log"
	// "go-game/mapgen/loader"
)

func main() {
	// mg := loader.MapLoader{}
	mg, err := infinitegen.New(10, 30)
	if err != nil {
		log.Panicf("Map generation error: %s", err)
	}
	ge := engine.NewGameEngine(mg)
	window.Run(ge.GfxLoop, nil)
}
