package model

type PromptConfig struct {
	InitValue       string   `yaml:"value,omitempty" short:"p" usage:"The (initial) prompt to use"`
	InitAttachments []string `yaml:"attachments,omitempty" short:"a" usage:"The (initial) attachments to use"`
	MinRows         uint     `yaml:"min-rows,omitempty" usage:"The minimal number of rows the prompt should have"`
	MaxRows         uint     `yaml:"max-rows,omitempty" usage:"The maximal number of rows the prompt should have"`
	PinTop          bool     `yaml:"pin-top,omitempty" usage:"Pin the prompt input at the top of the window (otherwise pin at the bottom)"`
	SubmitShortcut  Shortcut `yaml:"submit,omitempty" usage:"The shortcut for submit the prompt: "`
}

func (p *PromptConfig) SetDefaults() {
	p.MinRows = 1
	p.MaxRows = 4
	p.SubmitShortcut = Shortcut{Binding: []string{"Alt+Enter", "Alt+NumpadEnter"}}
	p.PinTop = true
}

func (p *PromptConfig) Validate() error {
	if ve := p.SubmitShortcut.Validate(); ve != nil {
		return ve
	}

	return nil
}
