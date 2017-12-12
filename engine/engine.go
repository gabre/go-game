package engine

import (
	"fmt"
	"go-game/mapgen"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"azul3d.org/engine/lmath"
	"azul3d.org/engine/mouse"
	"log"
)

const camZoomSpeed = 0.01 // 0.01x zoom for each scroll wheel click.
const camMinZoom = 0.1

type GameEngine struct {
	mapGenerator mapgen.MapGenerator
	cam          *camera.Camera
	d            gfx.Device
	camZoom      float64
	w            window.Window
}

func NewGameEngine(generator mapgen.MapGenerator) *GameEngine {
	return &GameEngine{mapGenerator: generator, camZoom: 1.0}
}

func (e *GameEngine) updateCamera() {
	if e.camZoom < camMinZoom {
		e.camZoom = camMinZoom
	}
	setOrthoScale(e.cam, e.d.Bounds(), e.camZoom)
}

// gfxLoop is responsible for drawing things to the window.
func (e *GameEngine) GfxLoop(w window.Window, d gfx.Device) {
	e.w = w
	e.d = d
	// Create a new orthographic (2D) camera.
	e.cam = camera.NewOrtho(d.Bounds())

	// Update the camera now.
	e.updateCamera()

	// Move the camera back two units away from the card.
	e.cam.SetPos(lmath.Vec3{0, -2, 0})

	// Load TMX map file.
	log.Printf("GfxLoop\n")
	layers, err := e.mapGenerator.GenerateMap(0, 0)
	log.Printf("Layers: %s\n", layers)

	if err != nil {
		log.Panicf("Error while generating map: ", err)
	}

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, getEventMask())

	for {
		log.Printf("Start draw\n")
		// Handle events.
		window.Poll(events, e.handlEvent)

		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Draw the TMX map to the screen.
		for _, layer := range layers {
			for _, obj := range layer {
				d.Draw(d.Bounds(), obj, e.cam)
			}
		}

		// Render the whole frame.
		d.Render()
	}
}

func (e *GameEngine) handlEvent(event window.Event) {
	switch ev := event.(type) {
	case window.FramebufferResized:
		// Update the camera's to account for the new width and height.
		e.updateCamera()

	case mouse.ButtonEvent:
		if ev.Button == mouse.Left && ev.State == mouse.Up {
			// Toggle mouse grab.
			props := e.w.Props()
			props.SetCursorGrabbed(!props.CursorGrabbed())
			e.w.Request(props)
		}

	case mouse.Scrolled:
		// Zoom and update the camera.
		e.camZoom -= ev.Y * camZoomSpeed
		e.updateCamera()

	case window.CursorMoved:
		if ev.Delta {
			p := lmath.Vec3{ev.X, 0, -ev.Y}
			p = p.MulScalar(e.camZoom)
			e.cam.SetPos(e.cam.Pos().Add(p))
		}

	case keyboard.Typed:
		switch ev.S {
		case "m":
			// Toggle MSAA now.
			msaa := !e.d.MSAA()
			e.d.SetMSAA(msaa)
			fmt.Println("MSAA Enabled?", msaa)
		case "r":
			e.cam.SetPos(lmath.Vec3{0, -2, 0})
		}
	}
}
