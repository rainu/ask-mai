package config

type PromptConfig struct {
	InitValue       string   `config:"value" short:"p" usage:"The (initial) prompt to use"`
	InitAttachments []string `config:"attachments" short:"a" usage:"The (initial) attachments to use"`
	MinRows         uint     `config:"min-rows" usage:"The minimal number of rows the prompt should have"`
	MaxRows         uint     `config:"max-rows" usage:"The maximal number of rows the prompt should have"`
	SubmitShortcut  Shortcut `config:"submit" usage:"The shortcut for submit the prompt: "`
}
