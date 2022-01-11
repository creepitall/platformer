package character

import "github.com/faiface/pixel"

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

func (p *Physics) Validate() bool {
	if p.Rectangle.Min.X < 0 {
		return false
	}
	// Больше размера сцены
	// TODO здесь должна быть константа
	if p.Rectangle.Max.X > 2880 {
		return false
	}

	return true
}

func (p *Physics) Update(dt float64, ctrl pixel.Vec) {
	//if p.Validate() {
		p.updateSideX(dt, &ctrl)
	//}
}

// Обновить физические данные движения по X
func (p *Physics) updateSideX(dt float64, ctrl *pixel.Vec) {
	switch {
	case ctrl.X < 0:
		p.Velocity.X = -p.RunSpeed
	case ctrl.X > 0:
		p.Velocity.X = +p.RunSpeed
	default:
		p.Velocity.X = 0
	}

	p.Rectangle = p.Rectangle.Moved(p.Velocity.Scaled(dt))
}

func (p *Physics) ReturnRectangleW() float64 {
	return p.Rectangle.W()
}

func (p *Physics) ReturnRectangleH() float64 {
	return p.Rectangle.H()
}

func (p *Physics) ReturnRectangleCenter() pixel.Vec {
	return p.Rectangle.Center()
}
