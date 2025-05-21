package model

import "github.com/rainu/ask-mai/internal/config/model/llm"

type Profile struct {
	Meta ProfileMeta `yaml:",inline,omitempty"`

	Printer PrinterConfig `yaml:"print,omitempty"`
	LLM     llm.LLMConfig `yaml:"llm,omitempty"`
	UI      UIConfig      `yaml:"ui,omitempty"`
	History History       `yaml:"history,omitempty"`

	RestartShortcut Shortcut `yaml:"restart-shortcut,omitempty" usage:"The shortcut for triggering a restart: "`
}

type ProfileMeta struct {
	Icon        string `yaml:"icon,omitempty" usage:"Profile icon"`
	Description string `yaml:"description,omitempty" usage:"Profile description"`
}

func (c *Profile) SetDefaults() {
	if c.RestartShortcut.Binding == nil {
		c.RestartShortcut = Shortcut{Binding: []string{"Alt+KeyR"}}
	}
}

func (c *Profile) Validate() error {
	if ve := c.RestartShortcut.Validate(); ve != nil {
		return ve
	}

	if ve := c.UI.Validate(); ve != nil {
		return ve
	}

	if ve := c.LLM.Validate(); ve != nil {
		return ve
	}

	if ve := c.Printer.Validate(); ve != nil {
		return ve
	}

	if ve := c.History.Validate(); ve != nil {
		return ve
	}

	return nil
}
