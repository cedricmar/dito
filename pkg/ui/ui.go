package ui

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/cedricmar/editor/pkg/viewmodel"
)

type Screen struct {
	view *viewmodel.View
	buf  *bytes.Buffer
}

func NewUI() *Screen {
	return &Screen{
		buf: bytes.NewBuffer([]byte{}),
	}
}

func (s *Screen) ListenUpdates(upChan <-chan *viewmodel.View) {
	go func() {
		for updatedView := range upChan {
			s.view = updatedView
			s.Paint()
		}
	}()
}

func (s *Screen) Paint() {
	s.buf.Reset()
	s.showCursor(false)
	s.moveCursor(0, 0)

	s.drawRows()

	s.moveCursor(s.view.Cursor.X, s.view.Cursor.Y)
	s.showCursor(true)

	s.buf.WriteTo(os.Stdout)
}

func (s *Screen) ClearScreen() {
	s.buf.Reset()
	s.moveCursor(0, 0)
	s.buf.WriteString("\x1b[2J") // Empty screen
	s.buf.WriteTo(os.Stdout)
}

func (s *Screen) drawRows() {
	for y := 1; y < s.view.Viewport.Height; y++ {
		s.clearLine()
		if y < s.view.Viewport.Height-1 {
			s.buf.WriteString("~")
			if y == s.view.Viewport.Height/3 {
				s.drawWelcome()
			}
			s.buf.WriteString("\r\n")
		} else {
			s.drawStatusBar()
		}
	}
}

func (s *Screen) drawWelcome() {
	w := fmt.Sprintf("Dito editor -- version %s", "0.0.1")
	wlen := len(w)
	if wlen > s.view.Viewport.Width-1 {
		wlen = s.view.Viewport.Width - 1
	}
	padLen := (s.view.Viewport.Width - wlen) / 2
	pad := ""
	if padLen > 0 {
		pad = strings.Repeat(" ", padLen-1)
	}

	b := make([]byte, len(pad)+wlen)
	copy(b, pad+w)
	s.buf.Write(b)
}

func (s *Screen) clearLine() {
	s.buf.WriteString("\x1b[2K")
}

func (s *Screen) clearLineFromCursor() {
	s.buf.WriteString("\x1b[K")
}

func (s *Screen) clearLineBeforeCursor() {
	s.buf.WriteString("\x1b[1K")
}

func (s *Screen) drawStatusBar() {
	ib := s.view.Infobar
	fmt.Fprintf(s.buf, "%s | %d:%d %d:%d | debug: %s",
		ib.Status,
		s.view.Cursor.X,
		s.view.Cursor.Y,
		s.view.Viewport.Width,
		s.view.Viewport.Height,
		s.view.Infobar.Debug,
	)
}

func (s *Screen) moveCursor(x, y int) {
	fmt.Fprintf(s.buf, "\x1b[%d;%dH", y+1, x+1)
}

func (s *Screen) showCursor(show bool) {
	if show {
		s.buf.WriteString("\x1b[?25h")
		return
	}
	s.buf.WriteString("\x1b[?25l")

}
