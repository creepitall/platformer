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
)

type heroPhys struct {
	gravity   float64
	runSpeed  float64
	jumpSpeed float64

	rect   pixel.Rect
	vel    pixel.Vec
	ground bool
}

type heroAnim struct {
	sheet pixel.Picture
	anims []pixel.Rect
	rate  float64

	//state   animState
	counter float64
	dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
}

var ObjFrames []pixel.Rect

var CurrentHeroPhysics *heroPhys

var CurrentHeroAnimation *heroAnim

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

// for test
func initHeroPlayer(assetsHero pixel.Picture) {
	CurrentHeroPhysics = &heroPhys{
		gravity:   -512,
		runSpeed:  96,
		jumpSpeed: 192,
		rect:      pixel.R(32, 64, 96, 128),
	}

	CurrentHeroAnimation = &heroAnim{
		sheet: assetsHero,
		anims: returHeroRect(assetsHero),
		rate:  1.0 / 10,
		dir:   +1,
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 960, 480),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	backgroundSprites := returnBackGroundSprite()
	assets, frontSprites := returnFrontSprite()

	rock1 := pixel.NewSprite(assets, frontSprites[6])
	rock2 := pixel.NewSprite(assets, frontSprites[7])

	assetsHero, err := loadPicture("assets/KnightRun_scale.png")
	if err != nil {
		panic(err)
	}
	//hero := pixel.NewSprite(assetsHero, assetsHero.Bounds())


	initHeroPlayer(assetsHero)


	camPos := pixel.ZV

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// cam := pixel.IM.Moved(pixel.V(0, 0))
		// win.SetMatrix(cam)

		camPos = pixel.Lerp(camPos, pixel.ZV, 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		win.SetMatrix(cam)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			fmt.Println(win.MousePosition())
		}

		if win.JustPressed(pixelgl.KeyR) {
			initHeroPlayer(assetsHero)
		}

		ctrl := pixel.ZV
		if win.Pressed(pixelgl.KeyLeft) {
			//camPos.X -= camSpeed * dt
			ctrl.X--
		}
		if win.Pressed(pixelgl.KeyRight) {
			//camPos.X += camSpeed * dt
			ctrl.X++
		}
		if win.JustPressed(pixelgl.KeyUp) {
			ctrl.Y = 1
		}

		win.Clear(color.White)

		//sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(win.Bounds().Center()))
		//fmt.Println(win.Bounds().Center())
		for _, sprite := range backgroundSprites {
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(480, 240)))
		}
		// for _, sprite := range backgroundSprites {
		// 	sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(480*3, 240)))
		// }

		rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(32, 32)))
		rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(96, 32)))
		rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(160, 32)))
		rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(224, 32)))

		CurrentHeroPhysics.update(dt, ctrl)
		CurrentHeroAnimation.update(dt, CurrentHeroPhysics)

		//hero.Set(ga.sheet, ga.frame)
		//hero.Draw(win, pixel.IM.Scaled(pixel.ZV, 3.0).Moved(hp.rect.Center()))

		CurrentHeroAnimation.draw(win, CurrentHeroPhysics)

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

func returnBackGroundSprite() []*pixel.Sprite {
	var background pixel.Picture
	var sprite []*pixel.Sprite
	var err error

	var spritesName = []string{
		"assets/background1.png",
		"assets/background3.png",
		"assets/background4b.png",
	}

	for _, name := range spritesName {
		background, err = loadPicture(name)
		if err != nil {
			panic(err)
		}
		sprite = append(sprite, pixel.NewSprite(background, background.Bounds()))
	}

	return sprite
}

func returnFrontSprite() (pixel.Picture, []pixel.Rect) {
	var assets pixel.Picture
	var ObjFrames []pixel.Rect

	assets, err := loadPicture("assets/build_3.png")
	if err != nil {
		panic(err)
	}

	for y := 0.0; y < assets.Bounds().Max.Y; y += 32.0 {
		for x := 0.0; x < assets.Bounds().Max.X; x += 32.0 {
			ObjFrames = append(ObjFrames, pixel.R(x, y, x+32.0, y+32.0))
		}
	}

	return assets, ObjFrames
}

func returHeroRect(assets pixel.Picture) []pixel.Rect {
	var anims []pixel.Rect
	for x := 0.0; x < assets.Bounds().Max.X; x += 32.0 {
		anims = append(anims, pixel.R(x, 0, x+32.0, 32))
	}

	return anims
}

func (hp *heroPhys) update(dt float64, ctrl pixel.Vec) {
	// apply controls
	switch {
	case ctrl.X < 0:
		hp.vel.X = -hp.runSpeed
	case ctrl.X > 0:
		hp.vel.X = +hp.runSpeed
	default:
		hp.vel.X = 0
	}
	// platform --
	//hp.ground = false

	hp.vel.Y += hp.gravity * dt
	hp.rect = hp.rect.Moved(hp.vel.Scaled(dt))

	hp.ground = false
	if hp.vel.Y <= 0 {
		if ((hp.rect.Max.X + hp.rect.Min.X) / 2) <= 265 {
			hp.ground = true
		}

		if hp.ground {
			if hp.rect.Max.Y < 128 {
				hp.rect = hp.rect.Moved(pixel.V(0, 64-hp.rect.Min.Y))
				hp.vel.Y = 0
			}
		}
	}


	if ctrl.Y > 0 {
		hp.vel.Y = hp.jumpSpeed
	}

}

func (ha *heroAnim) update(dt float64, phys *heroPhys) {
	ha.counter += dt

	var state string = "staying"
	if phys.vel.Len() > 0 {
		state = "running"
	}
	//if !phys.ground {
	//	state = "jumping"
	//}

	switch state {
	case "staying":
		ha.frame = ha.anims[0]
	case "running":
		i := int(math.Floor(ha.counter / ha.rate))
		ha.frame = ha.anims[i%len(ha.anims)]
	case "jumping":

	}

	if phys.vel.X != 0 {
		if phys.vel.X > 0 {
			ha.dir = +1
		} else {
			ha.dir = -1
		}
	}
}

func (ha *heroAnim) draw(t pixel.Target, phys *heroPhys) {
	if ha.sprite == nil {
		ha.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// // draw the correct frame with the correct position and direction
	ha.sprite.Set(ha.sheet, ha.frame)
	//ha.sprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 3.0).Moved(phys.rect.Center()))
	ha.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.rect.W()/ha.sprite.Frame().W(),
			phys.rect.H()/ha.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(+ha.dir, 1)).
		Moved(phys.rect.Center()),
	)
}
