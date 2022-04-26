package store

import (
	"github.com/cedricmar/editor/pkg/viewmodel"
)

// Store manages the state.

// There could be several stores.
// A store listens for all actions and decides on which of them to act.

// It then runs the action and emits another event for the UI to update.

type Store struct {
	// previousStates []state
	state *state
}

type Runner interface {
	Run()
}

type state struct {
	v             int
	cursor        cursor
	winSize       size
	status, flash string
	debug         string
}

type cursor struct {
	row, col int
	show     bool
}

type size struct {
	rows, cols int
}

var storeInst *Store

func CreateStore() chan Runner {
	storeInst = &Store{
		state: &state{
			v:       1,
			cursor:  cursor{row: 0, col: 0, show: true},
			winSize: size{rows: 0, cols: 0},
			status:  "<C-x> to exit",
			flash:   "Welcome to Dito | <C-x> to exit",
			debug:   "",
		},
	}

	return make(chan Runner)
}

func ListenActions(actions <-chan Runner) chan *viewmodel.View {
	updateChan := make(chan *viewmodel.View)
	go func() {
		for action := range actions {
			action.Run()
			mv := ModelView()
			updateChan <- &mv
		}
	}()
	return updateChan
}

func getState() *state {
	return storeInst.state
}

func getCursor() *cursor {
	return &getState().cursor
}

func getWinSize() *size {
	return &getState().winSize
}
