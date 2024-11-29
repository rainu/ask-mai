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

	tabs            *container.AppTabs
	mdOutput        *widget.RichText
	mdOutputScroll  *container.Scroll
	rawOutput       *widget.Entry
	rawOutputScroll *container.Scroll
	inText          *components.TextInput

	btnStop      *widget.Button
	btnClipboard *widget.Button
}

func NewWindow(controller *Controller, prompt string, width, height uint) *Window {
	r := &Window{}

	r.app = app.New()
	r.window = r.app.NewWindow("Input - Ask mAI")
	r.window.Resize(fyne.NewSize(float32(width), 0))
	r.window.CenterOnScreen()

	r.mdOutput = widget.NewRichTextFromMarkdown("")
	r.mdOutput.Wrapping = fyne.TextWrapWord

	r.mdOutputScroll = container.NewScroll(r.mdOutput)
	r.mdOutputScroll.SetMinSize(fyne.NewSize(float32(width), float32(height)))
	r.mdOutputScroll.Hide()

	r.rawOutput = widget.NewMultiLineEntry()
	r.rawOutput.Wrapping = fyne.TextWrapWord

	r.rawOutputScroll = container.NewScroll(r.rawOutput)
	r.rawOutputScroll.SetMinSize(fyne.NewSize(float32(width), float32(height)))
	r.rawOutputScroll.Hide()

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

	r.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Rendered", theme.DocumentIcon(), r.mdOutputScroll),
		container.NewTabItemWithIcon("Plain", theme.DocumentPrintIcon(), r.rawOutputScroll),
	)
	r.tabs.Hide()

	if controller.backendBuilder.Type == backend.TypeSingleShot {
		r.window.SetContent(container.NewBorder(inputContainer, r.btnClipboard, nil, nil, r.tabs))
	} else if controller.backendBuilder.Type == backend.TypeMultiShot {
		r.window.SetContent(container.NewBorder(nil, inputContainer, nil, nil, r.tabs))
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
