package error

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
	app    fyne.App
	window fyne.Window
}

func NewWindow(errorText string) *Window {
	r := &Window{}

	r.app = app.New()
	r.window = r.app.NewWindow("Error - Ask mAI")
	r.window.SetContent(container.NewVBox(widget.NewLabel(errorText)))

	if deskCanvas, ok := r.window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyEscape {
				r.app.Quit()
			}
		})
	}

	return r
}

func (w *Window) ShowAndRun() {
	w.window.ShowAndRun()
}
