package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	camSpeed     = 500.0
	camZoom      = 1.5
	camZoomSpeed = 1.5
	// rocks    []*pixel.Sprite
	trees []*pixel.Sprite
	//	matrices []pixel.Matrix
	CurrentSprite spriteSettings
)

var ObjFrames []pixel.Rect

type spriteSettings struct {
	blockSize float64
	scale     float64
}

type spritesheet struct {
	ss pixel.Picture
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	CurrentSprite = spriteSettings{
		blockSize: 32.0,
		scale:     1.0,
	}

	assets, err := loadPicture("assets/build_3.png")
	if err != nil {
		panic(err)
	}

	ss := spritesheet{
		ss: assets,
	}
	ss.initSprites()
	//objects, matrices := ss.createLevel()

	// spritesheet, err := loadPicture("assets/trees.png")
	// if err != nil {
	// 	panic(err)
	// }

	//tree := pixel.NewSprite(spritesheet, pixel.R(0, 0, 32, 32))

	background, err := loadPicture("assets/background3.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(background, background.Bounds())

	//black := color.RGBA{23, 18, 23, 255}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	fmt.Println(pixel.V(1.0, 1.0))

	scene := pixel.R(30.0, 624.0, 990, 144.0)
	result := createGrid(scene)
	for idY, valueY := range result {
		for idX, valueX := range valueY {
			fmt.Printf("Y: %v X: %v coor: %v \r\n", idY, idX, valueX)
		}
	}

	camPos := scene.Center()

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		//camPos = pixel.Lerp(camPos, pixel.R(0, 0, 6, 7).Center(), 1-math.Pow(1.0/128, dt))
		//cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		//cam := pixel.IM.Moved(camPos.Scaled(-1))

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			// objects = append(objects, tree)
			// matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()))
			fmt.Println(win.MousePosition())
		}

		if win.Pressed(pixelgl.KeyLeft) {
			//if camPos.X > 0 {
			camPos.X -= camSpeed * dt
			//}
		}
		if win.Pressed(pixelgl.KeyRight) {
			//if camPos.X < 960 {
			camPos.X += camSpeed * dt
			//}
		}
		if win.Pressed(pixelgl.KeyDown) {
			//if camPos.Y > 0 {
			camPos.Y -= camSpeed * dt
			//}
		}
		if win.Pressed(pixelgl.KeyUp) {
			//if camPos.Y < 480 {
			camPos.Y += camSpeed * dt
			//}
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(color.White)

		sceneCenter := scene.Center()
		for i := 1.0; i < 4.0; i++ {
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(sceneCenter.X*i, sceneCenter.Y)))
		}
		//sprite.Draw(win, pixel.IM.Moved(scene.Center()))
		//sprite.Draw(win, pixel.IM.Moved(pixel.ZV))

		obj := pixel.NewSprite(ss.ss, returnFrame(1))
		//obj.Draw(win, pixel.IM.Moved(pixel.ZV))

		////obj.Draw(win, pixel.IM.Moved(scene.Min))

		//createFrame(ss.ss, scene, win)
		createMapNew(result, obj, win)

		// for i, obj := range objects {
		// 	obj.Draw(win, matrices[i])
		// }

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

func createFrame(p pixel.Picture, scene pixel.Rect, win *pixelgl.Window) {
	min := scene.Min
	max := scene.Max

	obj := pixel.NewSprite(p, returnFrame(1))

	blockSize := CurrentSprite.blockSize * CurrentSprite.scale

	for y := min.Y; y > max.Y; y -= blockSize {
		for x := min.X; x < max.X; x += blockSize {
			obj.Draw(win, pixel.IM.Moved(pixel.V(x+blockSize, y-blockSize)))
		}
	}

}

func createMapNew(currentScene [][]pixel.Vec, obj *pixel.Sprite, win *pixelgl.Window) {
	var currentMap = [15][30]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 2, 1, 2, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{2, 1, 2, 1, 2, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 9, 9, 9, 9, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 1, 2, 1, 2, 2},
		{9, 9, 9, 3, 9, 9, 3, 9, 2, 1, 2, 1, 2, 1, 2, 2, 2, 1, 2, 1, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0},
		{9, 3, 4, 9, 4, 9, 9, 9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 4, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0},
	}

	for lineId, line := range currentMap {
		for blockId, block := range line {
			if block != 0 {
				obj.Draw(win, pixel.IM.Scaled(pixel.ZV, CurrentSprite.scale).Moved(currentScene[lineId][blockId]))
			}
		}
	}
}

func createGrid(scene pixel.Rect) [][]pixel.Vec {
	min := scene.Min
	max := scene.Max

	blockSize := CurrentSprite.blockSize * CurrentSprite.scale

	var currentScene [][]pixel.Vec
	for y := min.Y; y > max.Y; y -= blockSize {
		var sceneXLine []pixel.Vec
		for x := min.X; x < max.X; x += blockSize {
			sceneXLine = append(sceneXLine, pixel.V(x+blockSize, y-blockSize))
		}
		currentScene = append(currentScene, sceneXLine)
	}

	return currentScene
}

func (s *spritesheet) createLevel() (objects []*pixel.Sprite, matrices []pixel.Matrix) {
	// x=32, y=16
	var currentMap = [16][32]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 2, 1, 2, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{2, 1, 2, 1, 2, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 9, 9, 9, 9, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 1, 2, 1, 2, 2, 2, 1},
		{9, 9, 9, 3, 9, 9, 3, 9, 2, 1, 2, 1, 2, 1, 2, 2, 2, 1, 2, 1, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 3, 4, 9, 4, 9, 9, 9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 4, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	blockSize := (CurrentSprite.blockSize * CurrentSprite.scale)
	maxY := blockSize * 12

	delta := pixel.V(1, maxY) // 32x32 sprite

	objects = make([]*pixel.Sprite, 0)
	matrices = make([]pixel.Matrix, 0)

	for _, line := range currentMap {
		for blockId, block := range line {
			a := blockSize * float64(blockId)
			if block > 0 {
				obj := pixel.NewSprite(s.ss, returnFrame(block))
				objects = append(objects, obj)
				newDelta := pixel.V(a+blockSize, delta.Y)
				matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, CurrentSprite.scale).Moved(newDelta))
			}
		}
		delta = pixel.V(1, delta.Y-(blockSize))
	}

	return objects, matrices
}

func (s *spritesheet) initSprites() {
	for y := 0.0; y < s.ss.Bounds().Max.Y; y += 32.0 {
		for x := 0.0; x < s.ss.Bounds().Max.X; x += 32.0 {
			ObjFrames = append(ObjFrames, pixel.R(x, y, x+32.0, y+32.0))
		}
	}
}

func returnFrame(vl int) pixel.Rect {
	var returnValue pixel.Rect
	switch vl {
	case 1: // up one
		returnValue = ObjFrames[6]
	case 2: // up two
		returnValue = ObjFrames[7]
	case 3: // under one
		returnValue = ObjFrames[1]
	case 4: // under two
		returnValue = ObjFrames[2]
	case 5: // up left
		returnValue = ObjFrames[3]
	case 6: // up right
		returnValue = ObjFrames[5]
	case 7: // dark one
		returnValue = ObjFrames[0]
	case 9: // dark two
		returnValue = ObjFrames[4]
	}

	return returnValue
}
