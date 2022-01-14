package character

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type physics interface {
	Update(dt float64, ctrl *pixel.Vec, cs State, platform []pixel.Rect)
	ReturnRectangleSumX() float64
	ReturnRectangleSumY() float64
	ReturnRectangleW() float64
	ReturnRectangleH() float64
	ReturnRectangleCenter() pixel.Vec
	ReturnInfo() string
}

type animation interface {
	Update(dt float64, ctrl pixel.Vec, cs CharacterState)
	ReturnFrameW() float64
	ReturnFrameH() float64
	Draw(t *pixelgl.Window, scaleXYVec pixel.Vec, rectCenter pixel.Vec)
	ReturnInfo() string
}

type state interface {
	Update(characterState CharacterState)
	ReturnCurrentState() CharacterState
	ReturnState() State
}

// Персонаж
type Character struct {
	IsPlayer  bool
	Physics   physics
	Animation animation
	State     state
}

func CreateNewCharacter(isPlayer bool, ph physics, an animation, st state) *Character {
	return &Character{
		IsPlayer:  isPlayer,
		Physics:   ph,
		Animation: an,
		State:     st,
	}
}

func (c *Character) Update(windows *pixelgl.Window, dt float64, ctrl pixel.Vec, platform []pixel.Rect) {
	var (
		//velocity float64
	)

	c.Physics.Update(dt, &ctrl, c.State.ReturnState(), platform)

	//c.State.Update(velocity)

	c.Animation.Update(dt, ctrl, c.State.ReturnCurrentState())

	scaleXYVec := pixel.V(
		c.Physics.ReturnRectangleW()/c.Animation.ReturnFrameW(),
		c.Physics.ReturnRectangleH()/c.Animation.ReturnFrameH(),
	)
	rectCentr := c.Physics.ReturnRectangleCenter()
	c.Animation.Draw(windows, scaleXYVec, rectCentr)
}
