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

func (t *TextInput) KeyUp(key *fyne.KeyEvent) {
	t.Entry.KeyUp(key)
	if t.OnKeyUp != nil {
		t.OnKeyUp(key)
	}
}
