package engine

import (
	"engo.io/ecs"
	"engo.io/engo/common"
	"engo.io/engo"
	"engo.io/engo/math"
)

var zeroX, zeroY = float32(0), float32(0)

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
// TODO this should be a method of engine
func (r *ChunkSystem) Update(dt float32) {
	defer r.engine.renderNearestMapPartsIfNeeded()
	UpdateViewPointWithCamera(r.engine, r.cameraSys)
}
func UpdateViewPointWithCamera(engine *GameEngine, cameraSys *common.CameraSystem) {
	halfWW := engo.WindowWidth() / 2
	halfWH := engo.WindowHeight() / 2
	// zero is (zeroX - halfW), so camera X measured from this is camX
	camX := cameraSys.X() + (zeroX - halfWW)
	camY := cameraSys.Y() + (zeroY - halfWH)
	actualViewPointX := camScalarToDiscretePoint(camX, engine.chunkWidth)
	actualViewPointY := camScalarToDiscretePoint(camY, engine.chunkHeight)
	engine.viewPoint = point{actualViewPointX, actualViewPointY}
}

func camScalarToDiscretePoint(camScalar float32, unitLength float32) int64 {
	translation := unitLength / 2
	camScalar += translation
	units := camScalar / unitLength
	return int64(math.Floor(units))
}