package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"azul3d.org/engine/lmath"
	"azul3d.org/engine/mouse"
	"azul3d.org/engine/tmx"
)

type MapGenerator interface {
	GenerateMap() (*tmx.Map, map[string]map[string]*gfx.Object)
}

type MapLoader struct {}

func (m MapLoader) GenerateMap() (*tmx.Map, map[string]map[string]*gfx.Object) {
	tmxMap, layers, err := tmx.LoadFile(defaultMapFile, nil)
	if err != nil {
		log.Panicf("Error when loading map: %s", err)
	}
	return tmxMap, layers
}

type GameEngine struct {
	mapGenerator MapGenerator
	cam *camera.Camera
	d gfx.Device
	camZoom float64
	w window.Window
}

func NewGameEngine(generator MapGenerator) *GameEngine {
	return &GameEngine{mapGenerator: generator, camZoom: 1.0}
}

// setOrthoScale sets the camera's projection matrix to an orthographic one
// using the given viewing rectangle. It performs scaling with the viewing
// rectangle.
func setOrthoScale(c *camera.Camera, view image.Rectangle, scale float64) {
	w := float64(c.View.Dx())
	w *= scale
	w = float64(int((w / 2.0)))

	h := float64(c.View.Dy())
	h *= scale
	h = float64(int((h / 2.0)))

	m := lmath.Mat4Ortho(-w, w, -h, h, c.Near, c.Far)
	c.P = gfx.ConvertMat4(m)
}

func (e *GameEngine) updateCamera() {
	if e.camZoom < camMinZoom {
		e.camZoom = camMinZoom
	}
	setOrthoScale(e.cam, e.d.Bounds(), e.camZoom)
}

func getEventMask() window.EventMask {
	// Create an event mask for the events we are interested in.
	evMask := window.FramebufferResizedEvents
	evMask |= window.CursorMovedEvents
	evMask |= window.MouseEvents
	evMask |= window.MouseScrolledEvents
	evMask |= window.KeyboardTypedEvents
	return evMask
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
	tmxMap, layers := e.mapGenerator.GenerateMap()

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, getEventMask())

	for {
		// Handle events.
		window.Poll(events, e.handlEvent)

		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Draw the TMX map to the screen.
		for _, layer := range tmxMap.Layers {
			objects, ok := layers[layer.Name]
			if ok {
				for _, obj := range objects {
					d.Draw(d.Bounds(), obj, e.cam)
				}
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


const defaultMapFile = "data/test_base64.tmx"
const camZoomSpeed = 0.01 // 0.01x zoom for each scroll wheel click.
const camMinZoom = 0.1

func main() {
	mg := MapLoader{}
	ge := NewGameEngine(mg)
	window.Run(ge.GfxLoop, nil)
}