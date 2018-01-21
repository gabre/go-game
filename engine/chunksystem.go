package engine

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"engo.io/engo"
	// "fmt"
	"fmt"
)

type ChunkSystem struct {
	engine *GameEngine
	cameraSys *common.CameraSystem
	r int64 // number of chunks rendered vertically/horizontally counted from the central chunk
			// e.g. [ -r ... CENTRAL ... +r ]
}

func NewChunkSystem(engine *GameEngine, cameraSys *common.CameraSystem, r int64) *ChunkSystem {
	return &ChunkSystem{engine, cameraSys, r}
}

// Remove is called whenever an Entity is removed from the World, in order to remove it from this sytem as well
func (*ChunkSystem) Remove(ecs.BasicEntity) {}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (r *ChunkSystem) Update(dt float32) {
	halfWW := engo.WindowWidth() / 2
	halfWH := engo.WindowHeight() / 2
	camX := r.cameraSys.X() + (halfWW)
	camY := r.cameraSys.Y() + (halfWH)

	//camX := r.cameraSys.X()
	//camY := r.cameraSys.Y()

	top    := r.engine.mapTopLeftCoord.Y
	bottom := r.engine.mapTopLeftCoord.Y + r.engine.height
	left   := r.engine.mapTopLeftCoord.X
	right  := r.engine.mapTopLeftCoord.X + r.engine.width
//float32(0) //
	l := float32(0) // float32(math.Max(float64(r.engine.width / 3), float64(r.engine.height / 3)))
	topLimit    := top + l
	bottomLimit := bottom - l
	leftLimit   := left + l
	rightLimit  := right - l

	defer r.engine.renderNearestMapPartsIfNeeded()

	//
	// 1 | o | 2
	// o | x | o
	// 3 | o | 4
	//

	fmt.Printf("T B L R %f %f %f %f   L X Y %f %f %f  %d %d\n", top, bottom, left, right, l, camX, camY, r.engine.viewPoint.x, r.engine.viewPoint.y)
	if (topLimit > camY && leftLimit > camX) {
		println(1)
		r.engine.viewPoint.x -= 1
		r.engine.viewPoint.y -= 1
		r.engine.mapTopLeftCoord.X -= r.engine.width
		r.engine.mapTopLeftCoord.Y -= r.engine.height
		return
	}
	if (topLimit > camY && rightLimit < camX) {
		println(2)
		r.engine.viewPoint.x += 1
		r.engine.viewPoint.y -= 1
		r.engine.mapTopLeftCoord.X += r.engine.width
		r.engine.mapTopLeftCoord.Y -= r.engine.height
		return
	}
	if (bottomLimit < camY && leftLimit > camX) {
		println(3)
		r.engine.viewPoint.x -= 1
		r.engine.viewPoint.y += 1
		r.engine.mapTopLeftCoord.X -= r.engine.width
		r.engine.mapTopLeftCoord.Y += r.engine.height
		return
	}
	if (bottomLimit < camY && rightLimit < camX) {
		println(4)
		r.engine.viewPoint.x += 1
		r.engine.viewPoint.y += 1
		r.engine.mapTopLeftCoord.X += r.engine.width
		r.engine.mapTopLeftCoord.Y += r.engine.height
		return
	}

	//
	// o | 1 | o
	// 2 | x | 3
	// o | 4 | o
	//

	if (topLimit > camY) {
		println(5)
		r.engine.viewPoint.y -= 1
		r.engine.mapTopLeftCoord.Y -= r.engine.height
		return
	}
	if (leftLimit > camX) {
		println(6)
		r.engine.viewPoint.x -= 1
		r.engine.mapTopLeftCoord.X -= r.engine.width
		return
	}
	if (rightLimit < camX) {
		println(7)
		r.engine.viewPoint.x += 1
		r.engine.mapTopLeftCoord.X += r.engine.width
		return
	}
	if (bottomLimit < camY) {
		println(8)
		r.engine.viewPoint.y += 1
		r.engine.mapTopLeftCoord.Y += r.engine.height
		return
	}

	//if (r.cameraSys.X() < r.engine.viewPoint.X || r.cameraSys.Y() < r.engine.viewPoint.Y ||
	//	(r.engine.viewPoint.X + sx) < r.cameraSys.X() || (r.engine.viewPoint.Y + sy) < r.cameraSys.Y())
}