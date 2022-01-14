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

func (s *State) Update(characterState CharacterState) {
	s.CurrentState = characterState
}

func (c *CharacterState) FromString(value string) CharacterState {
	switch value{
	case "Jump":
		return CharStateJump
	case "Stay":
		return CharStateStay
	case "Run":
		return CharStateRun
	default:
		return CharStateStay
	}
}

func (c *CharacterState) ToString(value CharacterState) string {
	switch value {
	case CharStateJump:
		return "Jump"
	case CharStateStay:
		return "Stay"
	case CharStateRun:
		return "Run"
	default:
		return ""
	}
}

func (s *State) ReturnCurrentState() CharacterState {
	return s.CurrentState
}

func (s *State) ReturnState() State {
	return State{CurrentState: s.CurrentState}
}