package scene

import (
	"fmt"
	"github.com/creepitall/platformer/internal/character"
	"github.com/creepitall/platformer/internal/domain"
	"github.com/creepitall/platformer/internal/image"
	"github.com/creepitall/platformer/internal/pkg/config"
	"github.com/creepitall/platformer/internal/worldmap"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"log"
	"time"
)

var (
	frames = 0
	second = time.Tick(time.Second)
	scene  Scene
	char   *character.Character
	//
	ctrl pixel.Vec
	camPos pixel.Vec
)

type Scene struct {
	Init         bool   // Сцена загружена
	CurrentScene string // Имя сцены
	SceneSprites map[string][]*pixel.Sprite // Спрайты для сцены
	PlatformPhysics []pixel.Rect
	Test1        *pixel.Sprite              // Передний план сцены (временно так)
}

func ConsDataToLog() string {
	var str string
	str += char.Physics.ReturnInfo()
	str += "\r\n"
	str += char.Animation.ReturnInfo()

	return str
}

func TikTak(dt float64, isLog bool) {
	if !isLog {
		return
	}

	exptTime := domain.PreviousTime.Add(2 * time.Second)

	if time.Until(exptTime) < (0 * time.Second) {
		chp := ConsDataToLog()
		log.Println(chp)
		log.Printf("dt: %v \r\n", dt)
		domain.PreviousTime = time.Now()
	}
}

func DrawScene(windows *pixelgl.Window, config *config.Config, last *time.Time) {
	if !scene.Init {
		InitCharacter()
		scene.InitScene()
		scene.InitSceneObjects()
	}

	// Ограничение FPS
	if config.EnableFPS {
		time.Sleep(1 * time.Second / 60)
	}

	// Коэф.
	dt := time.Since(*last).Seconds()
	*last = time.Now()

	// Камера
	cameraMatrix := ReturnMatrix(windows, dt)
	windows.SetMatrix(cameraMatrix)

	ctrl = pixel.ZV
	checkIO(windows)

	windows.Clear(color.White)

	scene.Draw(windows)
	char.Update(windows, dt, ctrl, scene.PlatformPhysics)

	windows.Update()

	TikTak(dt, config.Logging)

	// Вывод данных FPS
	frames++
	select {
	case <-second:
		windows.SetTitle(fmt.Sprintf("%s | FPS: %d", config.Title, frames))
		frames = 0
	default:
	}

}

// Инициализация персонажа
func InitCharacter() {
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

func (s *Scene) InitSceneObjects() {
	value := worldmap.CreateNewMap()

	s.PlatformPhysics = make([]pixel.Rect, 0, len(value))
	for _, vl := range value {
		rc := pixel.Rect{
			Min: pixel.V(vl.Min.X, vl.Min.Y),
			Max: pixel.V(vl.Max.X, vl.Max.Y),
		}

		s.PlatformPhysics = append(s.PlatformPhysics, rc)
	}
}

// Отрисовка сцены
func (s *Scene) Draw(w pixel.Target) {
	for _, sprite := range s.SceneSprites["back"] {
		sprite.Draw(w, pixel.IM.Scaled(pixel.ZV, 4.0).Moved(pixel.V(1280, 960)))
	}

	s.Test1.Draw(w, pixel.IM.Scaled(pixel.ZV, 1.0).Moved(pixel.V(1440, 672)))
}
