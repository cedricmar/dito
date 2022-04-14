package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"strconv"
	"strings"

	"github.com/cedricmar/editor/pkg/key"
	"github.com/cedricmar/editor/pkg/term"
	"github.com/cedricmar/editor/pkg/ui"
)

func main() {

	t, err := term.GetTermState()
	if err != nil {
		fmt.Printf("term: %+v %v\n", t, err)
		os.Exit(1)
	}

	err = t.RawMode()
	if err != nil {
		fmt.Printf("term set RAW_MODE: %v\n", err)
		os.Exit(1)
	}

	scr := ui.NewUI(t.TerminalHeight(), t.TerminalWidth())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan bool, 1)

	for {
		scr.Paint()
		key.HandleKeypress(stop)

		select {
		default:
		case <-stop:
			goto FINISH
		case <-sigs:
			goto FINISH
		}
	}

FINISH:
	scr.ClearScreen()
	err = t.Reset()
	if err != nil {
		fmt.Printf("term set SANE_MODE: %v\n", err)
	}

	os.Exit(0)

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
