package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"time"

	"github.com/creepitall/test_pixel/internal/image"
	"github.com/creepitall/test_pixel/internal/models"

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
	isDeath   bool
	isJump    bool

	rect   pixel.Rect
	vel    pixel.Vec
	ground bool
}

type heroAnim struct {
	sheet pixel.Picture
	//	anims []pixel.Rect
	rate float64

	//state   animState
	counter float64
	dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
}

var ObjFrames []pixel.Rect

var CurrentHeroPhysics *heroPhys

var CurrentHeroAnimation *heroAnim

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "little story: the knight",
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

	rock1 := models.SceneSprites["front"][6]
	rock2 := models.SceneSprites["front"][7]

	///
	// basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// basicTxt := text.New(pixel.V(800, 470), basicAtlas)

	// fmt.Fprintf(basicTxt, "dead status: %v \r\n", CurrentHeroPhysics.isDeath)
	// fmt.Fprintf(basicTxt, "jump status: %v", CurrentHeroPhysics.isJump)
	///

	camPos := pixel.ZV

	last := time.Now()
	for !win.Closed() {
		time.Sleep(1 * time.Second / 60) // fix to 60 fps
		dt := time.Since(last).Seconds()
		last = time.Now()

		camPos = pixel.Lerp(camPos, pixel.ZV, 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		win.SetMatrix(cam)

		if CurrentHeroPhysics.isDeath {
			//time.Sleep(1 * time.Second)
			initHeroPlayer()
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			fmt.Println(win.MousePosition())
		}

		if win.JustPressed(pixelgl.KeyR) {
			initHeroPlayer()
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

		for _, sprite := range models.SceneSprites["back"] {
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(480, 240)))
		}

		rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(32, 32)))
		rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(96, 32)))
		rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(160, 32)))
		rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(224, 32)))

		CurrentHeroPhysics.update(dt, ctrl)
		CurrentHeroAnimation.update(dt, CurrentHeroPhysics)

		//hero.Set(ga.sheet, ga.frame)
		//hero.Draw(win, pixel.IM.Scaled(pixel.ZV, 3.0).Moved(hp.rect.Center()))

		CurrentHeroAnimation.draw(win, CurrentHeroPhysics)

		// basicTxt.Color = colornames.Whitesmoke
		// basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 1))

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
	initGameConfig()
	initHeroPlayer()
	pixelgl.Run(run)
}

func initGameConfig() {
	models.CurrentScene = "start"

	image.FillFrontSpriteByScene()
	image.FillHeroPlayerSprite()
}

func initHeroPlayer() {
	CurrentHeroPhysics = &heroPhys{
		gravity:   -512,
		runSpeed:  96,
		jumpSpeed: 192,
		rect:      pixel.R(32, 64, 96, 128),
		isDeath:   false,
		isJump:    false,
	}

	CurrentHeroAnimation = &heroAnim{
		sheet: models.HeroPlayerStayAssets,
		//	anims: returHeroRect(assetsHero),
		rate: 1.0 / 10,
		dir:  +1,
	}
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
		if ((hp.rect.Max.X+hp.rect.Min.X)/2) <= 265 && hp.rect.Min.Y >= 32 {
			hp.ground = true
		}

		if hp.ground {
			if hp.rect.Max.Y < 128 {
				hp.rect = hp.rect.Moved(pixel.V(0, 64-hp.rect.Min.Y))
				hp.vel.Y = 0
				hp.isJump = false
			}
		}
	}

	if !hp.isJump && ctrl.Y > 0 {
		hp.vel.Y = hp.jumpSpeed
		hp.isJump = true
	}

	if hp.rect.Max.Y < 0 {
		hp.isDeath = true
	}

}

func (ha *heroAnim) update(dt float64, phys *heroPhys) {
	ha.counter += dt

	var state string = "staying"
	if phys.vel.Len() > 0 {
		state = "running"
	}
	if phys.isJump {
		state = "jumping"
	}

	//i := int(math.Floor(ha.counter / ha.rate))
	switch state {
	case "staying":
		i := int(math.Floor(ha.counter / ha.rate))
		ha.sheet = models.HeroPlayerStayAssets
		ha.frame = models.HeroPlayerStayFrames[i%len(models.HeroPlayerStayFrames)]
	case "running":
		i := int(math.Floor(ha.counter / ha.rate))
		ha.sheet = models.HeroPlayerRunAssets
		ha.frame = models.HeroPlayerRunFrames[i%len(models.HeroPlayerRunFrames)]
	case "jumping":
		ha.sheet = models.HeroPlayerJumpAssets
		speed := phys.vel.Y
		i := int((-speed/phys.jumpSpeed + 1) / 2 * float64(len(models.HeroPlayerJumpFrames)))
		if i < 0 {
			i = 0
		}
		if i >= len(models.HeroPlayerJumpFrames) {
			i = len(models.HeroPlayerJumpFrames) - 1
		}
		ha.frame = models.HeroPlayerJumpFrames[i]
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
	// draw the correct frame with the correct position and direction
	ha.sprite.Set(ha.sheet, ha.frame)
	ha.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.rect.W()/ha.sprite.Frame().W(),
			phys.rect.H()/ha.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(+ha.dir, 1)).
		Moved(phys.rect.Center()),
	)
}
