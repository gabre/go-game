package engine

import (
	"image"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/lmath"
)

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

func getEventMask() window.EventMask {
	// Create an event mask for the events we are interested in.
	evMask := window.FramebufferResizedEvents
	evMask |= window.CursorMovedEvents
	evMask |= window.MouseEvents
	evMask |= window.MouseScrolledEvents
	evMask |= window.KeyboardTypedEvents
	return evMask
}
