package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TextInput struct {
	widget.Entry

	OnKeyUp func(key *fyne.KeyEvent)
}

func NewTextInput() *TextInput {
	result := &TextInput{}
	result.ExtendBaseWidget(result)

	return result
}

func (e *TextInput) KeyUp(key *fyne.KeyEvent) {
	e.Entry.KeyUp(key)
	if e.OnKeyUp != nil {
		e.OnKeyUp(key)
	}
}
