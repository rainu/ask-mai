package config

type PromptConfig struct {
	InitValue       string   `yaml:"value" short:"p" usage:"The (initial) prompt to use"`
	InitAttachments []string `yaml:"attachments" short:"a" usage:"The (initial) attachments to use"`
	MinRows         uint     `yaml:"min-rows" usage:"The minimal number of rows the prompt should have"`
	MaxRows         uint     `yaml:"max-rows" usage:"The maximal number of rows the prompt should have"`
	PinTop          bool     `yaml:"pin-top" usage:"Pin the prompt input at the top of the window (otherwise pin at the bottom)"`
	SubmitShortcut  Shortcut `yaml:"submit" usage:"The shortcut for submit the prompt: "`
}
