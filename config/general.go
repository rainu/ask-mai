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

type DebugConfig struct {
	LogLevel               int    `yaml:"log-level"`
	PprofAddress           string `yaml:"pprof-address" usage:"Address for the pprof server (only available for debug binary)"`
	OpenInspectorOnStartup bool   `yaml:"open-inspector-on-startup" usage:"Open the inspector on startup (only available for debug binary)"`
	PrintVersion           bool   `config:"version" yaml:"-" short:"v" usage:"Show the version"`
}

func (c *Config) GetUsage(field string) string {
	switch field {
	case "LogLevel":
		return fmt.Sprintf("Log level (debug(%d), info(%d), warn(%d), error(%d))", slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError)
	}
	return ""
}

func (c *Config) Validate() error {
	if c.Debug.LogLevel < int(slog.LevelDebug) || c.Debug.LogLevel > int(slog.LevelError) {
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
