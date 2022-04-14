package term

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type TermState struct {
	InitialTermios *unix.Termios
	*unix.Termios
	*unix.Winsize
}

func GetTermState() (*TermState, error) {
	stdinFD := int(os.Stdin.Fd())
	t, err := getTermios(stdinFD)
	if err != nil {
		return nil, err
	}
	w, err := getWinsize(stdinFD)
	if err != nil {
		return nil, err
	}
	return &TermState{
		InitialTermios: tcopy(*t),
		Termios:        t,
		Winsize:        w,
	}, nil
}

func (t *TermState) SetTermios() error {
	stdinFD := int(os.Stdin.Fd())
	return unix.IoctlSetTermios(stdinFD, uint(unix.TIOCSETA), t.Termios)
}

func (t *TermState) Reset() error {
	t.Termios = t.InitialTermios
	return t.SetTermios()
}

func (t *TermState) RawMode() error {
	t.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	t.Oflag &^= unix.OPOST
	t.Cflag |= unix.CS8
	t.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	t.Cc[unix.VMIN] = 0
	t.Cc[unix.VTIME] = 1
	return t.SetTermios()
}

func (t *TermState) TerminalWidth() int {
	return int(t.Winsize.Col)
}

func (t *TermState) TerminalHeight() int {
	return int(t.Winsize.Row)
}

func (t *TermState) Print() {
	fmt.Printf("Termios: %+v - Winsize: %+v\r\n", t.Termios, t.Winsize)
}

func getTermios(fd int) (*unix.Termios, error) {
	return unix.IoctlGetTermios(fd, uint(unix.TIOCGETA))
}

func getWinsize(fd int) (*unix.Winsize, error) {
	// @TODO - write a fallback in case IoctlGetWinsize fails
	return unix.IoctlGetWinsize(fd, uint(unix.TIOCGWINSZ))
}

func tcopy(t unix.Termios) *unix.Termios {
	return &t
}
