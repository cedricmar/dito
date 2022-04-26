package keyboard

import (
	"os"

	"github.com/cedricmar/editor/pkg/action"
	"github.com/cedricmar/editor/pkg/store"
)

const (
	EscSeqChar = '\x1b'
)

func Isctrl(b byte) bool {
	return b <= 31 || b == 127
}

func Ctrl(key byte) byte {
	return key & 0x1f
}

func HandleKeypress(actionCh chan<- store.Runner, stop chan bool) {
	c := readKey()

	switch c {
	case Ctrl('j'):
		actionCh <- action.CursorMoveAction(action.ActionCursorMoveLeft)
	case Ctrl('i'):
		actionCh <- action.CursorMoveAction(action.ActionCursorMoveUp)
	case Ctrl('k'):
		actionCh <- action.CursorMoveAction(action.ActionCursorMoveDown)
	case Ctrl('l'):
		actionCh <- action.CursorMoveAction(action.ActionCursorMoveRight)
	case Ctrl('x'):
		stop <- true
	}
}

func readKey() byte {
	var b []byte = make([]byte, 1)
	for {
		nread, err := os.Stdin.Read(b)
		if nread != 1 || err != nil {
			return EscSeqChar
		}

		// Handle multibyte inputs
		if b[0] == EscSeqChar {
			seq := make([]byte, 3)

			nc := make([]byte, 1)
			nread, err := os.Stdin.Read(nc)
			if nread != 1 || err != nil {
				return EscSeqChar
			}
			seq[0] = nc[0]

			nread, err = os.Stdin.Read(nc)
			if nread != 1 || err != nil {
				return EscSeqChar
			}
			seq[1] = nc[0]

			if seq[0] == '[' {
				switch seq[1] {
				case 'A': // Up
					return Ctrl('i')
				case 'B': // Down
					return Ctrl('k')
				case 'C': // Right
					return Ctrl('l')
				case 'D': // Left
					return Ctrl('j')
				}
			}

			return EscSeqChar
		} else {
			return b[0]
		}
	}
}
