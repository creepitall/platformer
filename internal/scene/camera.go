package scene

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type Camera pixel.Vec

var (
	cam Camera
)

// TODO: здесь не нужен pixelgl
func ReturnMatrix(windows *pixelgl.Window, dt float64) pixel.Matrix {
	cam.Change(windows)

	camPos = pixel.Lerp(camPos, cam.ReturnVec(), 1-math.Pow(1.0/128, dt))

	if camPos.Y < 0 {
		camPos.Y = 0
	} else if camPos.Y >= windows.Bounds().H()/2-100 {
		camPos.Y = windows.Bounds().H()/2 - 100
	}
	if camPos.X < 0 {
		camPos.X = 0
	} else if camPos.X >= (windows.Bounds().W() / 2) {
		camPos.X = (windows.Bounds().W() / 2)
	}

	return pixel.IM.Moved(camPos.Scaled(-1))
}

func (c *Camera) ReturnVec() pixel.Vec {
	return pixel.V(c.X, c.Y)
}

func (c *Camera) Change(windows *pixelgl.Window) {
	c.X = char.Physics.ReturnRectangleSumX() / 2
	c.Y = char.Physics.ReturnRectangleSumY() / 2
	c.X -= (windows.Bounds().W() / 2)
	c.Y -= (windows.Bounds().H() / 2)
}
