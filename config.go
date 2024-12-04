package main

import (
	"flag"
	"fmt"
	"github.com/rainu/ask-mai/backend/copilot"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	BackendCopilot     = "copilot"
	BackendOpenAI      = "openai"
	BackendAnythingLLM = "anythingllm"

	PrinterFormatPlain = "plain"
	PrinterFormatJSON  = "json"
	PrinterTargetOut   = "stdout"
	PrinterTargetErr   = "stderr"
)

type Config struct {
	Prompt string
	Width  uint
	Height uint

	Backend     string
	OpenAI      OpenAIConfig
	AnythingLLM AnythingLLMConfig

	Printer PrinterConfig

	LogLevel int
}

type OpenAIConfig struct {
	APIKey       string
	SystemPrompt string
}

type AnythingLLMConfig struct {
	BaseURL   string
	Token     string
	Workspace string
}

type PrinterConfig struct {
	Format  string
	Targets []io.WriteCloser
	targets string
}

func ParseConfig() *Config {
	c := &Config{}

	flag.IntVar(&c.LogLevel, "ll", int(slog.LevelError), fmt.Sprintf("Log level (debug(%d), info(%d), warn(%d), error(%d))", slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError))
	flag.UintVar(&c.Width, "width", 1920/2, "The (initial) width of the window")
	flag.UintVar(&c.Height, "height", 1080/3, "The maximal height of the chat response area")
	flag.StringVar(&c.Prompt, "prompt", "", "The prompt to use")
	flag.StringVar(&c.Backend, "backend", BackendCopilot, fmt.Sprintf("The backend to use ('%s', '%s', '%s')", BackendCopilot, BackendOpenAI, BackendAnythingLLM))

	flag.StringVar(&c.OpenAI.APIKey, "openai-api-key", "", "OpenAI API Key")
	flag.StringVar(&c.OpenAI.SystemPrompt, "openai-system-prompt", "", "OpenAI System Prompt")

	flag.StringVar(&c.AnythingLLM.BaseURL, "anythingllm-base-url", "", "Base URL for AnythingLLM")
	flag.StringVar(&c.AnythingLLM.Token, "anythingllm-token", "", "Token for AnythingLLM")
	flag.StringVar(&c.AnythingLLM.Workspace, "anythingllm-workspace", "", "Workspace for AnythingLLM")

	flag.StringVar(&c.Printer.Format, "print-format", PrinterFormatJSON, fmt.Sprintf("Response printer format (%s, %s)", PrinterFormatPlain, PrinterFormatJSON))
	flag.StringVar(&c.Printer.targets, "print-targets", PrinterTargetOut, fmt.Sprintf("Comma seperated response printer targets (%s, %s, <path/to/file>)", PrinterTargetOut, PrinterTargetErr))

	flag.Parse()

	for _, target := range strings.Split(c.Printer.targets, ",") {
		target = strings.TrimSpace(target)

		if target == PrinterTargetOut {
			c.Printer.Targets = append(c.Printer.Targets, os.Stdout)
		} else if target == PrinterTargetErr {
			c.Printer.Targets = append(c.Printer.Targets, os.Stderr)
		} else {
			file, err := os.Create(target)
			if err != nil {
				panic(fmt.Errorf("Error creating printer target file: %w", err))
			}
			c.Printer.Targets = append(c.Printer.Targets, file)
		}
	}

	return c
}

func (c Config) Validate() error {
	if c.LogLevel < int(slog.LevelDebug) || c.LogLevel > int(slog.LevelError) {
		return fmt.Errorf("Invalid log level")
	}

	if c.Backend != BackendCopilot && c.Backend != BackendOpenAI && c.Backend != BackendAnythingLLM {
		return fmt.Errorf("Invalid backend")
	}

	if c.Backend == BackendCopilot && !copilot.IsCopilotInstalled() {
		return fmt.Errorf("GitHub Copilot is not installed")
	}

	if c.Backend == BackendOpenAI && c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API Key is missing")
	}

	if c.Backend == BackendAnythingLLM {
		if c.AnythingLLM.BaseURL == "" {
			return fmt.Errorf("AnythingLLM Base URL is missing")
		}
		if c.AnythingLLM.Token == "" {
			return fmt.Errorf("AnythingLLM Token is missing")
		}
		if c.AnythingLLM.Workspace == "" {
			return fmt.Errorf("AnythingLLM Workspace is missing")
		}
	}

	if c.Printer.Format != PrinterFormatJSON && c.Printer.Format != PrinterFormatPlain {
		return fmt.Errorf("Invalid response printer format")
	}

	return nil
}
