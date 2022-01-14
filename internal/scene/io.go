package scene

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
)

func checkIO(windows *pixelgl.Window) {
	if windows.JustPressed(pixelgl.MouseButtonLeft) {
		fmt.Println(scene.Init)
	}

	if windows.JustPressed(pixelgl.KeyEscape) {
		windows.SetClosed(true)
	}

	if windows.Pressed(pixelgl.KeyLeft) {
		ctrl.X--
	}
	if windows.Pressed(pixelgl.KeyRight) {
		ctrl.X++
	}
	if windows.JustPressed(pixelgl.KeyUp) {
		ctrl.Y = 1
	}
}
