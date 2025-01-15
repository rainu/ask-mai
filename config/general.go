package config

import (
	"github.com/rainu/ask-mai/config/llm"
)

type Config struct {
	UI UIConfig `yaml:"ui"`

	LLM llm.LLMConfig `config:"" yaml:"llm"`

	Printer PrinterConfig `yaml:"print"`

	Debug DebugConfig `config:"" yaml:"debug"`

	Config string `config:"config" short:"c" yaml:"-" usage:"Path to the configuration yaml file"`
}

func (c *Config) Validate() error {
	if ve := c.Debug.Validate(); ve != nil {
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

	return nil
}
