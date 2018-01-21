package main

import (
	"go-game/mapgen/infinitegen"
	"log"
	"go-game/engine"
	"engo.io/engo"
	"path/filepath"
	"go-game/mapgen"
)

func main() {
	resolution := 3
	factory := func(res int) (mapgen.MapGenerator, error) {
		return infinitegen.New(10, res)
	}
	ge := engine.NewGameEngine(resolution, factory)
	opts := engo.RunOptions{
		Title:  "TileMap Demo",
		Width:  800,
		Height: 800,
		AssetsRoot: getCwd(),
	}
	engo.Run(opts, ge)
}

func getCwd() string {
	fp, err := filepath.Abs(".")
	if err != nil {
		log.Panicf("Unable to get current working directory.")
	}
	return fp
}
