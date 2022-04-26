package action

type ActionType int

const (
	ActionWinSize ActionType = 1000 + iota
	ActionCursorToggle
	ActionCursorMove
	ActionCursorMoveUp
	ActionCursorMoveDown
	ActionCursorMoveLeft
	ActionCursorMoveRight
	ActionDebugUpdate
)

type Action struct {
	Kind ActionType
	call func()
}

func NewAction(kind ActionType, caller func()) *Action {
	return &Action{
		Kind: kind,
		call: caller,
	}
}

func (a *Action) Run() {
	a.call()
}
