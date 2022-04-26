package term

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

	ts := &TermState{
		InitialTermios: tcopy(*t),
		Termios:        t,
	}

	ws, err := ts.getWinsize(stdinFD)
	if err != nil {
		return nil, err
	}
	ts.Winsize = ws

	return ts, nil
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

func (t *TermState) getWinsize(fd int) (*unix.Winsize, error) {
	ws, err := unix.IoctlGetWinsize(fd, uint(unix.TIOCGWINSZ))
	if ws == nil || err != nil {
		// Fallback in case IoctlGetWinsize fails
		t.RawMode()

		os.Stdout.Write([]byte("\x1b[999C\x1b[999B"))
		c, r, err := getCursorPosition()
		if err != nil {
			return nil, err
		}

		t.Reset()

		ws = &unix.Winsize{
			Row: uint16(r),
			Col: uint16(c),
		}
	}
	return ws, err
}

func getCursorPosition() (col, row int, err error) {
	os.Stdout.Write([]byte("\x1b[6n")) // Print cursor pos

	buf := make([]byte, 32)
	i := 0
	for i < len(buf)-1 {
		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil {
			return
		}
		buf[i] = b[0]
		if buf[i] == []byte("R")[0] {
			break
		}
		i++
	}
	buf[i] = []byte("\x00")[0]

	if buf[0] != []byte("\x1b")[0] || buf[1] != []byte("[")[0] {
		err = errors.New("wrong format")
		return
	}

	reader := bytes.NewReader(buf[2:])
	n, err := fmt.Fscanf(reader, "%d;%d", &row, &col)
	if n != 2 || err != nil {
		err = errors.New("wrong format")
	}
	return
}

func getTermios(fd int) (*unix.Termios, error) {
	return unix.IoctlGetTermios(fd, uint(unix.TIOCGETA))
}

func tcopy(t unix.Termios) *unix.Termios {
	return &t
}

// err := getSttyState(&originalSttyState)
// if err != nil {
// 	log.Fatal(err)
// }
// defer setSttyState(&originalSttyState)

// setSttyState(bytes.NewBufferString("cbreak"))
// setSttyState(bytes.NewBufferString("-echo"))

// var b []byte = make([]byte, 1)
// for {
// 	os.Stdin.Read(b)
// 	fmt.Printf("Read character: %s\n", b)
// }

/*
	err = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}*/

func getSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func progress(current, total, cols int) string {
	prefix := strconv.Itoa(current) + " / " + strconv.Itoa(total)
	bar_start := " ["
	bar_end := "] "

	bar_size := cols - len(prefix+bar_start+bar_end)
	amount := int(float32(current) / (float32(total) / float32(bar_size)))
	remain := bar_size - amount

	bar := strings.Repeat("X", amount) + strings.Repeat(" ", remain)
	return Bold(prefix) + bar_start + bar + bar_end
}
