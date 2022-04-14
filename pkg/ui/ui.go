package ui

import (
	"fmt"
	"os"
)

type coord struct {
	row, col int
}

type State struct {
	cursorPos coord
	winSize   coord
	status    string
}

func NewUI(r, c int) *State {
	return &State{
		cursorPos: coord{row: 0, col: 0},
		winSize:   coord{row: r, col: c},
		status:    "<C-x> to exit",
	}
}

func (s *State) Paint() {
	s.ClearScreen()
	s.drawRows()

	s.moveCursor(1, 1)
}

func (s *State) ClearScreen() {
	os.Stdout.Write([]byte("\x1b[2J")) // Empty screen
	s.moveCursor(1, 1)
}

func (s *State) moveCursor(row, col int) {
	cmd := fmt.Sprintf("\x1b[%d;%dH", row, col)
	os.Stdout.Write([]byte(cmd))
}

func (s *State) drawRows() {
	for y := 1; y < s.winSize.row; y++ {
		if y == s.winSize.row-1 {
			fmt.Fprintf(
				os.Stdout,
				" %d %d / %d %d | %s\r\n",
				s.cursorPos.col,
				s.cursorPos.row,
				s.winSize.col,
				s.winSize.row,
				s.status,
			)
			continue
		}
		os.Stdout.Write([]byte("~\r\n"))
	}
}
