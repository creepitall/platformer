package scene

import (
	"fmt"
	"github.com/creepitall/platformer/internal/character"
	"github.com/creepitall/platformer/internal/domain"
	"github.com/creepitall/platformer/internal/image"
	"github.com/creepitall/platformer/internal/pkg/config"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"time"
)

var (
	frames = 0
	second = time.Tick(time.Second)
	cam    pixel.Matrix
	scene  Scene
	char   *character.Character
	ctrl  pixel.Vec
)

type Scene struct {
	Init         bool   // Сцена загружена
	CurrentScene string // Имя сцены
	//DefaultSprites map[string]*pixel.Sprite
	SceneSprites map[string][]*pixel.Sprite // Спрайты для сцены
	Test1        *pixel.Sprite              // Передний план сцены (временно так)
}

func DrawScene(windows *pixelgl.Window, config *config.Config, last time.Time) {
	if !scene.Init {
		InitCharacter()
		scene.InitScene()
	}

	// Ограничение FPS
	if config.EnableFPS {
		time.Sleep(1 * time.Second / config.FPS)
	}

	// Коэф.
	last = time.Now()
	dt := time.Since(last).Seconds()

	// Камера
	cam = pixel.IM.Moved(pixel.ZV.Scaled(-1))
	windows.SetMatrix(cam)

	ctrl = pixel.ZV

	checkIO(windows)

	windows.Clear(color.White)

	scene.Draw(windows)
	char.Update(windows, dt, ctrl)

	windows.Update()

	// Вывод данных FPS
	frames++
	select {
	case <-second:
		windows.SetTitle(fmt.Sprintf("%s | FPS: %d", config.Title, frames))
		frames = 0
	default:
	}

}

func InitCharacter() {
	// Грузим персонажа
	ph := character.CreateNewPhysics(196, 300, pixel.R(32, 64, 96, 128), pixel.Vec{})

	image.FillHeroPlayerSprite()

	an := character.CreateNewAnimation(+1, domain.HeroPlayerStayAssets)
	st := character.CreateNewState()
	char = character.CreateNewCharacter(true, ph, an, st)
}

// Временно тут
// Инициализация сцены
func (s *Scene) InitScene() {
	// Грузим сцену
	scene.CurrentScene = "start"
	scene.SceneSprites = image.FillFrontSpriteByScene("start")
	scene.Test1 = scene.SceneSprites["front"][0]

	scene.Init = true
}

// Отрисовка сцены
func (s *Scene) Draw(w pixel.Target) {
	for _, sprite := range s.SceneSprites["back"] {
		sprite.Draw(w, pixel.IM.Scaled(pixel.ZV, 4.0).Moved(pixel.V(1280, 960)))
	}

	s.Test1.Draw(w, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(1440, 672)))
}
