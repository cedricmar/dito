package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cedricmar/editor/pkg/action"
	"github.com/cedricmar/editor/pkg/keyboard"
	"github.com/cedricmar/editor/pkg/store"
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

	actionChan := store.CreateStore()
	storeUpdateChan := store.ListenActions(actionChan)

	ui := ui.NewUI()
	ui.ListenUpdates(storeUpdateChan)

	actionChan <- action.WinSizeAction(t.TerminalWidth(), t.TerminalHeight())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	stop := make(chan bool, 1)

	for {
		keyboard.HandleKeypress(actionChan, stop)

		select {
		default:
		case <-stop:
			goto FINISH
		case <-sigs:
			goto FINISH
		}
	}

FINISH:
	ui.ClearScreen()
	err = t.Reset()
	if err != nil {
		fmt.Printf("term set SANE_MODE: %v\n", err)
	}

	os.Exit(0)
}
