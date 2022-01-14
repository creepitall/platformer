package character

type CharacterState string

const (
	CharStateJump  CharacterState = "Jump"
	CharStateStay  CharacterState = "Stay"
	CharStateRun   CharacterState = "Run"
	CharStateDeath CharacterState = "Death"
)

// Различные состояния персонажа
type State struct {
	CurrentState CharacterState
}

func CreateNewState() *State {
	return &State{}
}

func (s *State) Update(vel float64) {
	s.CurrentState = CharStateStay
	if vel != 0 {
		s.CurrentState = CharStateRun
	}
}

func (s *State) ReturnCurrentState() CharacterState {
	return s.CurrentState
}
