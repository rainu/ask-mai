package config

import (
	"fmt"
	"github.com/rainu/ask-mai/config/llm"
	"log/slog"
)

type Config struct {
	UI UIConfig `yaml:"ui"`

	LLM llm.LLMConfig `config:"" yaml:"llm"`

	Printer PrinterConfig `yaml:"print"`

	Debug DebugConfig `config:"" yaml:"debug"`
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
