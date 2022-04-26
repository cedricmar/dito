package action

import "github.com/cedricmar/editor/pkg/store"

// An action is an anonymous function that capture all
// information necessary to do an action.

// Actions have a type property identifying the action type,
// this is useful if you need to have several stores (switch).

func WinSizeAction(width, height int) *Action {
	return NewAction(ActionWinSize, func() {
		store.SetWinSize(height, width)
	})
}

func CursorSetPositionAction(x, y int) *Action {
	return NewAction(ActionCursorMove, func() {
		store.SetCursorPosition(x, y)
	})
}

func CursorMoveAction(at ActionType) *Action {
	return NewAction(at, func() {
		switch at {
		case ActionCursorMoveUp:
			store.SetCursorUp()
		case ActionCursorMoveDown:
			store.SetCursorDown()
		case ActionCursorMoveLeft:
			store.SetCursorLeft()
		case ActionCursorMoveRight:
			store.SetCursorRight()
		}
	})
}

func CursorToggleAction() *Action {
	return NewAction(ActionCursorToggle, func() {
		store.SetCursorVisibilitySwitch()
	})
}

func DebugSetAction(txt string) *Action {
	return NewAction(ActionDebugUpdate, func() {
		store.SetDebug(txt)
	})
}
