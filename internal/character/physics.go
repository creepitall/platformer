package character

import (
	"encoding/json"
	"github.com/creepitall/platformer/internal/domain"
	"github.com/faiface/pixel"
)

// Физика
type Physics struct {
	RunSpeed  float64    // Скорость бега
	JumpSpeed float64    // Сокорость прыжка
	Rectangle pixel.Rect // Физические границы персонажа
	Velocity  pixel.Vec  // Вектор скорости
}

func CreateNewPhysics(runSpeed, jumpSpeed float64, rec pixel.Rect, vel pixel.Vec) *Physics {
	return &Physics{
		RunSpeed:  runSpeed,
		JumpSpeed: jumpSpeed,
		Rectangle: rec,
		Velocity:  vel,
	}
}

func (p *Physics) ReturnInfo() string {
	type hp struct {
		RunSpeed  float64    `json:"runSpeed"`
		JumpSpeed float64    `json:"jumpSpeed"`
		Rectangle pixel.Rect `json:"rect"`
		Velocity  pixel.Vec  `json:"vel"`
	}
	tmp := &hp{
		RunSpeed:  p.RunSpeed,
		JumpSpeed: p.JumpSpeed,
		Rectangle: p.Rectangle,
		Velocity:  p.Velocity,
	}
	bytes, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func (p *Physics) Validate(ctrl *pixel.Vec) *pixel.Vec {
	// Движение <-
	if p.Rectangle.Min.X < 0 && ctrl.X < 0 {
		ctrl.X = 0
	}
	// Больше размера сцены, движение ->
	// TODO здесь должна быть константа
	if p.Rectangle.Max.X > 2880 && ctrl.X > 0 {
		ctrl.X = 0
	}

	return ctrl
}


func (p *Physics) Update(dt float64, ctrl *pixel.Vec, cs State, platform []pixel.Rect) {
	ctrl = p.Validate(ctrl)

	p.updateSideX(dt, ctrl, cs)
	p.updateSideY(dt, ctrl, cs, platform)
}

// Обновить физические данные движения по X
func (p *Physics) updateSideX(dt float64, ctrl *pixel.Vec, cs State) {
	switch {
	case ctrl.X < 0:
		p.Velocity.X = -p.RunSpeed
	case ctrl.X > 0:
		p.Velocity.X = +p.RunSpeed
	default:
		p.Velocity.X = 0
		cs.Update(CharStateRun)
	}

	p.Rectangle = p.Rectangle.Moved(p.Velocity.Scaled(dt))

	if p.Velocity.Len() > 0 {
		cs.Update(CharStateRun)
	}
}

func (p *Physics) updateSideY(dt float64, ctrl *pixel.Vec, cs State, platforms []pixel.Rect) {
	p.Velocity.Y += domain.GlobalGravity * dt

	if p.Velocity.Y <= 0 {
		avrX := p.ReturnRectangleSumX() / 2
		for _, platform := range platforms {
			if avrX <= platform.Min.X || avrX >= platform.Max.X {
				continue
			}
			if p.Rectangle.Min.Y >= platform.Min.Y {
				continue
			}
			p.Velocity.Y = 0
			p.Rectangle = p.Rectangle.Moved(pixel.V(0, platform.Min.Y-p.Rectangle.Min.Y))
			//hp.isJump = false
		}
	}
	if ctrl.Y > 0 {
		p.Velocity.Y = p.JumpSpeed
		//hp.isJump = true
		cs.Update(CharStateJump)
	}
}

func (p *Physics) ReturnRectangleW() float64 {
	return p.Rectangle.W()
}

func (p *Physics) ReturnRectangleH() float64 {
	return p.Rectangle.H()
}

func (p *Physics) ReturnRectangleSumX() float64 {
	return p.Rectangle.Min.X + p.Rectangle.Max.X
}

func (p *Physics) ReturnRectangleSumY() float64 {
	return p.Rectangle.Min.Y + p.Rectangle.Max.Y
}

func (p *Physics) ReturnRectangleCenter() pixel.Vec {
	return p.Rectangle.Center()
}
