package character

import (
	"encoding/json"
	"github.com/creepitall/platformer/internal/domain"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

// Анимация
type Animation struct {
	Sheet   pixel.Picture // Текущий видимый слой
	Rate    float64       // скорость смены фрейма
	Counter float64       // счетчик
	Dir     float64       // нужно для скалирования изображения по вертикали
	Frame   pixel.Rect    // Размер фрейма
	Sprite  *pixel.Sprite // Спрайты
}

func CreateNewAnimation(dir float64, sh pixel.Picture) *Animation {
	return &Animation{
		Sheet: sh,
		Rate:  1.0 / 10,
		Dir:   dir,
	}
}

func (a *Animation) ReturnInfo() string {
	type ha struct {
		Rate    float64    `json:"rate"`
		Counter float64    `json:"counter"`
		Dir     float64    `json:"dir"`
		Frame   pixel.Rect `json:"frame"`
	}
	tmp := &ha{
		Rate:    a.Rate,
		Counter: a.Counter,
		Dir:     a.Dir,
		Frame:   a.Frame,
	}
	bytes, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (a *Animation) Update(dt float64, ctrl pixel.Vec, cs CharacterState) {
	a.Counter += dt

	switch cs {
	case CharStateStay:
		i := int(math.Floor(a.Counter / a.Rate))
		a.Sheet = domain.HeroPlayerStayAssets
		a.Frame = domain.HeroPlayerStayFrames[i%len(domain.HeroPlayerStayFrames)]
	case CharStateRun:
		i := int(math.Floor(a.Counter / a.Rate))
		a.Sheet = domain.HeroPlayerRunAssets
		a.Frame = domain.HeroPlayerRunFrames[i%len(domain.HeroPlayerRunFrames)]
	}

	a.changeDir(ctrl)
}

func (a *Animation) changeDir(ctrl pixel.Vec) {
	if ctrl.X != 0 {
		if ctrl.X > 0 {
			a.Dir = +1
		} else {
			a.Dir = -1
		}
	}
}

// rectCenter - центр физической модели объекта
func (a *Animation) Draw(t *pixelgl.Window, scaleXYVec pixel.Vec, rectCenter pixel.Vec) {
	if a.Sprite == nil {
		a.Sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	a.Sprite.Set(a.Sheet, a.Frame)
	a.Sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, scaleXYVec).
		ScaledXY(pixel.ZV, pixel.V(+a.Dir, 1)).
		Moved(rectCenter),
	)
}

func (a *Animation) ReturnFrameW() float64 {
	return a.Frame.W()
}

func (a *Animation) ReturnFrameH() float64 {
	return a.Frame.H()
}
