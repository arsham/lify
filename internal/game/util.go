package game

import (
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

type gameContext string

const preLoadTimeStr gameContext = "preloadtime"

const (
	sceneLoading = "loading_scene"
	sceneLify    = "lify_scene"
)

// axis are the plural of axis.
type axis uint8

// This is an enum for what axes to centre around.
const (
	axixXY axis = iota
	axixX
	axixY
)

func putCentre(ctx *scene.Context, obj render.Renderable, ax axis) {
	objWidth, objHeight := obj.GetDims()
	wbds := ctx.Window.Bounds()
	switch ax {
	case axixXY:
		obj.SetPos(
			float64(wbds.X()/2-objWidth/2),
			float64(wbds.Y()-objHeight)/2,
		)
	case axixX:
		obj.SetPos(float64(wbds.X()-objWidth)/2, obj.Y())
	case axixY:
		obj.SetPos(obj.X(), float64(wbds.Y()-objHeight)/2)
	}
}
