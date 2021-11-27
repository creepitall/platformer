package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"time"

	"github.com/creepitall/test_pixel/internal/domain"
	"github.com/creepitall/test_pixel/internal/image"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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

	camera pixel.Vec
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

type screenLogger struct {
	bt *text.Text
	ba *text.Atlas
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title: "little story: the knight",
		//Bounds: pixel.R(0, 0, 960, 480),
		Bounds: pixel.R(0, 0, 1920, 960),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// rock1 := domain.SceneSprites["front"][6]
	// rock2 := domain.SceneSprites["front"][7]
	test1 := domain.SceneSprites["front"][0]

	camPos := pixel.ZV
	var cam pixel.Matrix

	// init text
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicScreenLogger := screenLogger{
		bt: text.New(pixel.V(10, 900), basicAtlas),
		ba: basicAtlas,
	}
	//

	last := time.Now()
	for !win.Closed() {
		//time.Sleep(1 * time.Second / 60) // fix to 60 fps
		dt := time.Since(last).Seconds()
		last = time.Now()

		//if CurrentHeroPhysics.rect.Max.X < (win.Bounds().Max.X / 2) {
		//	camPos = pixel.Lerp(camPos, pixel.ZV, 1-math.Pow(1.0/128, dt))
		//	//cam := pixel.IM.Moved(camPos.Scaled(-1))
		//	cam = pixel.IM.Moved(camPos.Scaled(-1))
		//} else {

		CurrentHeroPhysics.changeCameraValue(win)

		camPos = pixel.Lerp(camPos, CurrentHeroPhysics.camera, 1-math.Pow(1.0/128, dt))
		//camPos = pixel.Lerp(camPos, CurrentHeroPhysics.rect.Max, 1-math.Pow(1.0/128, dt))
		//cam = pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		if camPos.Y < 0 {
			camPos.Y = 0
		}
		// if camPos.Y > win.Bounds().H() {
		// 	camPos.Y = win.Bounds().H()
		// }
		if camPos.X < 0 {
			camPos.X = 0
		} else if camPos.X >= (win.Bounds().W() / 2) {
			camPos.X = (win.Bounds().W() / 2)
		}
		cam = pixel.IM.Moved(camPos.Scaled(-1))
		//}
		//cam := pixel.IM.Moved(CurrentHeroPhysics.rect.Center())
		win.SetMatrix(cam)

		if CurrentHeroPhysics.isDeath {
			//time.Sleep(1 * time.Second)
			initHeroPlayer()
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			//fmt.Println(win.MousePosition())
			fmt.Printf("char hero max[%v] \r\n", CurrentHeroPhysics.rect.Max)
			fmt.Println("")
			fmt.Printf("cam pos [%v] \r\n", camPos)
			fmt.Println("")
			fmt.Printf("current cp: [%v] \r\n zv: [%v] \r\n dt: %v \r\n values: %v \r\n", camPos, pixel.ZV, dt, 1-math.Pow(1.0/128, dt))
			fmt.Println("")

			// fmt.Printf("cam pos [%v] \r\n", camPos)
			// fmt.Printf("cam [%v] \r\n", cam)
		}

		if win.JustPressed(pixelgl.KeyR) {
			initHeroPlayer()
		}

		ctrl := pixel.ZV
		if win.Pressed(pixelgl.KeyLeft) {
			//camPos.X -= camSpeed * dt
			//camPos.X = CurrentHeroPhysics.rect.Max.X + 20
			ctrl.X--
		}
		if win.Pressed(pixelgl.KeyRight) {
			//
			//camPos.X = CurrentHeroPhysics.rect.Min.X - 20
			ctrl.X++
		}
		if win.JustPressed(pixelgl.KeyUp) {
			ctrl.Y = 1
		}

		win.Clear(color.White)

		for _, sprite := range domain.SceneSprites["back"] {
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(960, 480)))
		}

		// rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(32, 32)))
		// rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(96, 32)))
		// rock1.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(160, 32)))
		// rock2.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pixel.V(224, 32)))
		test1.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(1440, 672)))

		CurrentHeroPhysics.update(dt, ctrl)
		CurrentHeroAnimation.update(dt, CurrentHeroPhysics)

		//hero.Set(ga.sheet, ga.frame)
		//hero.Draw(win, pixel.IM.Scaled(pixel.ZV, 3.0).Moved(hp.rect.Center()))

		CurrentHeroAnimation.draw(win, CurrentHeroPhysics)

		basicScreenLogger.drawlog(win, cam)

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

func (sc *screenLogger) drawlog(win *pixelgl.Window, cam pixel.Matrix) {

	var position pixel.Vec = pixel.V(10, 900)
	if CurrentHeroPhysics.camera.X > -10.0 {
		position.X += CurrentHeroPhysics.camera.X
	}
	if CurrentHeroPhysics.camera.X >= (win.Bounds().W()/2)+10.0 {
		position.X = (win.Bounds().W() / 2) + 10
	}

	sc.bt = text.New(position, sc.ba)
	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintf(sc.bt, "text postion: %v \r\n", position)
	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintf(sc.bt, "dead status: %v \r\n", CurrentHeroPhysics.isDeath)
	fmt.Fprintf(sc.bt, "jump status: %v \r\n", CurrentHeroPhysics.isJump)

	fmt.Fprintf(sc.bt, "cam matrix: %v \r\n", cam)

	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintln(sc.bt, "current rect char:")
	sc.bt.Color = colornames.Red
	fmt.Fprintf(sc.bt, "[%.2f, %.2f, %.2f, %.2f] \r\n", CurrentHeroPhysics.rect.Min.X, CurrentHeroPhysics.rect.Min.Y, CurrentHeroPhysics.rect.Max.X, CurrentHeroPhysics.rect.Max.Y)

	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintln(sc.bt, "current camera char:")
	sc.bt.Color = colornames.Red
	fmt.Fprintf(sc.bt, "[%.2f, %.2f,] \r\n", CurrentHeroPhysics.camera.X, CurrentHeroPhysics.camera.Y)

	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintln(sc.bt, "current vel char:")
	sc.bt.Color = colornames.Red
	fmt.Fprintf(sc.bt, "[%.2f, %.2f] \r\n", CurrentHeroPhysics.vel.X, CurrentHeroPhysics.vel.Y)

	fmt.Fprintf(sc.bt, "windows W:%v, H:%v \r\n", win.Bounds().W(), win.Bounds().H())

	sc.bt.Draw(win, pixel.IM.Scaled(sc.bt.Orig, 1))

}

func main() {
	initGameConfig()
	initHeroPlayer()
	pixelgl.Run(run)
}

func initGameConfig() {
	domain.CurrentScene = "start"

	//worldmap.CreateNewMap()
	image.FillFrontSpriteByScene()
	image.FillHeroPlayerSprite()
}

func initHeroPlayer() {
	CurrentHeroPhysics = &heroPhys{
		gravity: -512,
		//runSpeed: 96,
		runSpeed:  196,
		jumpSpeed: 256,
		rect:      pixel.R(32, 64, 96, 128),
		camera:    pixel.ZV,
		isDeath:   false,
		isJump:    false,
	}

	CurrentHeroAnimation = &heroAnim{
		sheet: domain.HeroPlayerStayAssets,
		//	anims: returHeroRect(assetsHero),
		rate: 1.0 / 10,
		dir:  +1,
	}
}

func (hp *heroPhys) changeCameraValue(win *pixelgl.Window) {
	hp.camera.X = (hp.rect.Max.X + hp.rect.Min.X) / 2
	hp.camera.Y = (hp.rect.Max.Y + hp.rect.Min.Y) / 2
	hp.camera.X -= (win.Bounds().W() / 2)
	hp.camera.Y -= (win.Bounds().H() / 2)
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
		//if ((hp.rect.Max.X+hp.rect.Min.X)/2) <= 265 && hp.rect.Min.Y >= 32 {
		if ((hp.rect.Max.X+hp.rect.Min.X)/2) <= 2880 && hp.rect.Min.Y >= 32 {
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
		ha.sheet = domain.HeroPlayerStayAssets
		ha.frame = domain.HeroPlayerStayFrames[i%len(domain.HeroPlayerStayFrames)]
	case "running":
		i := int(math.Floor(ha.counter / ha.rate))
		ha.sheet = domain.HeroPlayerRunAssets
		ha.frame = domain.HeroPlayerRunFrames[i%len(domain.HeroPlayerRunFrames)]
	case "jumping":
		ha.sheet = domain.HeroPlayerJumpAssets
		speed := phys.vel.Y
		i := int((-speed/phys.jumpSpeed + 1) / 2 * float64(len(domain.HeroPlayerJumpFrames)))
		if i < 0 {
			i = 0
		}
		if i >= len(domain.HeroPlayerJumpFrames) {
			i = len(domain.HeroPlayerJumpFrames) - 1
		}
		ha.frame = domain.HeroPlayerJumpFrames[i]
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
