package question

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/ui/components"
)

type Window struct {
	app    fyne.App
	window fyne.Window

	mdOutput *widget.RichText
	inText   *components.TextInput

	btnStop      *widget.Button
	btnClipboard *widget.Button
}

func NewWindow(controller *Controller, prompt string) *Window {
	r := &Window{}

	r.app = app.New()
	r.window = r.app.NewWindow("Input - Ask mAI")
	r.window.Resize(fyne.NewSize(400, 0))
	r.window.CenterOnScreen()

	r.mdOutput = widget.NewRichTextFromMarkdown("")
	r.mdOutput.Wrapping = fyne.TextWrapWord
	r.mdOutput.Hide()

	r.inText = components.NewTextInput()
	r.inText.SetPlaceHolder("Enter question...")
	r.inText.OnSubmitted = controller.OnInputSubmitted
	r.inText.OnKeyUp = func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			r.app.Quit()
		}
	}

	r.btnStop = widget.NewButtonWithIcon("", theme.CancelIcon(), controller.OnClickStop)
	r.btnStop.Hide()

	r.btnClipboard = widget.NewButton("Copy", controller.OnClickClipboard)
	r.btnClipboard.Hide()

	inputContainer := container.NewBorder(nil, nil, nil, r.btnStop, r.inText)

	if controller.backendBuilder.Type == backend.TypeSingleShot {
		r.window.SetContent(container.NewBorder(inputContainer, r.btnClipboard, nil, nil, r.mdOutput))
	} else if controller.backendBuilder.Type == backend.TypeMultiShot {
		r.window.SetContent(container.NewBorder(nil, inputContainer, nil, nil, r.mdOutput))
	} else {
		panic("Unknown backend type!")
	}

	if deskCanvas, ok := r.window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyEscape {
				r.app.Quit()
			}
		})
	}

	controller.SetWindow(r)

	if prompt != "" {
		r.inText.SetText(prompt)
		controller.OnInputSubmitted(prompt)
	}

	return r
}

func (w *Window) ShowAndRun() {
	w.window.Canvas().Focus(w.inText)
	w.window.ShowAndRun()
}
