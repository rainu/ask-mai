package config

import (
	"fmt"
	"github.com/rainu/ask-mai/config/llm"
	"log/slog"
)

type Config struct {
	UI UIConfig `config:"ui"`

	LLM llm.LLMConfig `config:""`

	Printer PrinterConfig `config:"print"`

	LogLevel     int  `config:"log-level"`
	PrintVersion bool `config:"version" short:"v" usage:"Show the version"`
}

func (c *Config) GetUsage(field string) string {
	switch field {
	case "LogLevel":
		return fmt.Sprintf("Log level (debug(%d), info(%d), warn(%d), error(%d))", slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError)
	}
	return ""
}

func (c *Config) Validate() error {
	if c.LogLevel < int(slog.LevelDebug) || c.LogLevel > int(slog.LevelError) {
		return fmt.Errorf("Invalid log level")
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