package store

func SetWinSize(rows, cols int) {
	*getWinSize() = size{
		rows: rows,
		cols: cols,
	}
}

func SetCursorVisibilitySwitch() {
	getCursor().show = !getCursor().show
}

func SetCursorPosition(x, y int) {
	getCursor().col = x
	getCursor().row = y
}

func SetCursorUp() {
	if getCursor().row > 0 {
		getCursor().row--
	}
}

func SetCursorDown() {
	if getCursor().row < getWinSize().rows-1 {
		getCursor().row++
	}
}

func SetCursorLeft() {
	if getCursor().col > 0 {
		getCursor().col--
	}
}

func SetCursorRight() {
	if getCursor().col < getWinSize().cols-1 {
		getCursor().col++
	}
}

func SetDebug(txt string) {
	getState().debug = txt
}
