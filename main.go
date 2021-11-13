package main

import (
	"image"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	camPos       = pixel.ZV
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2
	// rocks    []*pixel.Sprite
	trees []*pixel.Sprite
	//	matrices []pixel.Matrix
	CurrentSprite spriteSettings
)

var ObjFrames []pixel.Rect

type spriteSettings struct {
	x     float64
	y     float64
	scale float64
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
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	CurrentSprite = spriteSettings{
		x:     32.0,
		y:     32.0,
		scale: 2.0,
	}

	assets, err := loadPicture("assets/build_3.png")
	if err != nil {
		panic(err)
	}

	ss := spritesheet{
		ss: assets,
	}
	ss.initSprites()
	objects, matrices := ss.createLevel()

	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	tree := pixel.NewSprite(spritesheet, pixel.R(0, 0, 32, 32))

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		//camPos = pixel.Lerp(camPos, pixel.R(0, 0, 6, 7).Center(), 1-math.Pow(1.0/128, dt))
		//cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		//cam := pixel.IM.Moved(camPos.Scaled(-1))

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			objects = append(objects, tree)
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()))
		}

		if win.Pressed(pixelgl.KeyLeft) {
			if camPos.X > 0 {
				camPos.X -= camSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			if camPos.Y > 0 {
				camPos.Y -= camSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Whitesmoke)

		for i, obj := range objects {
			obj.Draw(win, matrices[i])
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func (s *spritesheet) createLevel() (objects []*pixel.Sprite, matrices []pixel.Matrix) {
	// x=32, y=16
	var currentMap = [12][32]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 6, 1, 2, 1, 2, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{2, 1, 2, 1, 2, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 9, 9, 9, 9, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 3, 9, 9, 3, 9, 2, 1, 2, 1, 2, 1, 2, 2, 2, 1, 2, 1, 2, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 3, 4, 9, 4, 9, 9, 9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 3, 9, 9, 9, 9, 9, 9, 9, 4, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	blockSize := (CurrentSprite.x * CurrentSprite.scale)
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
