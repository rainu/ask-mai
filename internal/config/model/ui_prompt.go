package model

import "github.com/rainu/go-yacl"

type PromptConfig struct {
	MinRows        *uint    `yaml:"min-rows,omitempty" usage:"The minimal number of rows the prompt should have"`
	MaxRows        *uint    `yaml:"max-rows,omitempty" usage:"The maximal number of rows the prompt should have"`
	PinTop         *bool    `yaml:"pin-top,omitempty" usage:"Pin the prompt input at the top of the window (otherwise pin at the bottom)"`
	SubmitShortcut Shortcut `yaml:"submit,omitempty" usage:"The shortcut for submit the prompt: "`
}

func (p *PromptConfig) SetDefaults() {
	if p.MinRows == nil {
		p.MinRows = yacl.P(uint(1))
	}
	if p.MaxRows == nil {
		p.MaxRows = yacl.P(uint(4))
	}
	if p.SubmitShortcut.Binding == nil {
		p.SubmitShortcut = Shortcut{Binding: []string{"Alt+Enter", "Alt+NumpadEnter"}}
	}
	if p.PinTop == nil {
		p.PinTop = yacl.P(true)
	}
}

func (p *PromptConfig) Validate() error {
	if ve := p.SubmitShortcut.Validate(); ve != nil {
		return ve
	}

	return nil
}
