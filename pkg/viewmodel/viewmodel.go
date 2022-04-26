package viewmodel

type View struct {
	Cursor   Cursor
	Viewport Viewport
	Infobar  Infobar
}

type Cursor struct {
	X, Y int
	Show bool
}

type Viewport struct {
	Width, Height int
}

type Infobar struct {
	Status, Message, Debug string
}
