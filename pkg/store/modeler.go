package store

import (
	"fmt"

	vm "github.com/cedricmar/editor/pkg/viewmodel"
)

type ViewModeler interface {
	ModelView() vm.View
}

func ModelView() vm.View {
	getState().v += 1
	return vm.View{
		Cursor:   modelCursor(),
		Viewport: modelViewport(),
		Infobar:  modelInfobar(),
	}
}

func modelCursor() vm.Cursor {
	return vm.Cursor{
		X:    getState().cursor.col,
		Y:    getState().cursor.row,
		Show: getState().cursor.show,
	}
}

func modelViewport() vm.Viewport {
	return vm.Viewport{
		Width:  getState().winSize.cols,
		Height: getState().winSize.rows,
	}
}

func modelInfobar() vm.Infobar {
	return vm.Infobar{
		Status:  getState().status,
		Message: getState().flash,
		//Debug:   getState().debug,
		Debug: fmt.Sprintf("v: %d, m: %s", getState().v, getState().debug),
	}
}
