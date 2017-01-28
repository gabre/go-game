package main

import (
	"go-game/engine"
	"go-game/mapgen/loader"

	"azul3d.org/engine/gfx/window"
)

func main() {
	mg := loader.MapLoader{}
	ge := engine.NewGameEngine(mg)
	window.Run(ge.GfxLoop, nil)
}
