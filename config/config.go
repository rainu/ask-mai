package config

import (
	"flag"
	"fmt"
	"github.com/rainu/ask-mai/llms"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	BackendCopilot     = "copilot"
	BackendLocalAI     = "localai"
	BackendOpenAI      = "openai"
	BackendAnythingLLM = "anythingllm"
	BackendOllama      = "ollama"
	BackendMistral     = "mistral"
	BackendAnthropic   = "anthropic"

	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeSystem = "system"

	TranslucentNever = "never"
	TranslucentEver  = "ever"
	TranslucentHover = "hover"

	PrinterFormatPlain = "plain"
	PrinterFormatJSON  = "json"
	PrinterTargetOut   = "stdout"
	PrinterTargetErr   = "stderr"
)

type Config struct {
	UI UIConfig

	Backend     string
	LocalAI     LocalAIConfig
	OpenAI      OpenAIConfig
	AnythingLLM AnythingLLMConfig
	Ollama      OllamaConfig
	Mistral     MistralConfig
	Anthropic   AnthropicConfig
	CallOptions CallOptionsConfig

	Printer PrinterConfig

	LogLevel int
}

type UIConfig struct {
	Window       WindowConfig
	Prompt       string
	QuitShortcut Shortcut
	Theme        string
	CodeStyle    string
	Language     string
}

type WindowConfig struct {
	Title            string
	InitialWidth     string
	MaxHeight        string
	InitialPositionX string
	InitialPositionY string
	InitialZoom      float64
	BackgroundColor  WindowBackgroundColor
	StartState       int
	Frameless        bool
	Resizeable       bool
	Translucent      string
}

type WindowBackgroundColor struct {
	R uint
	G uint
	B uint
	A uint
}

type Shortcut struct {
	Code  string
	Alt   bool
	Ctrl  bool
	Meta  bool
	Shift bool
}

type PrinterConfig struct {
	Format  string
	Targets []io.WriteCloser
	targets string
}

func Parse(arguments []string) *Config {
	c := &Config{}

	flag.IntVar(&c.LogLevel, "ll", int(slog.LevelError), fmt.Sprintf("Log level (debug(%d), info(%d), warn(%d), error(%d))", slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError))

	flag.StringVar(&c.UI.Prompt, "ui-prompt", "", "The prompt to use")
	flag.StringVar(&c.UI.Window.Title, "ui-title", "Prompt - Ask mAI", "The window title")
	flag.StringVar(&c.UI.Window.InitialWidth, "ui-init-width", "CurrentScreen.Dimension.Width/2", "The (initial) width of the window")
	flag.StringVar(&c.UI.Window.MaxHeight, "ui-max-height", "CurrentScreen.Dimension.Height/3", "The maximal height of the chat response area")
	flag.StringVar(&c.UI.Window.InitialPositionX, "ui-init-pos-x", "CurrentScreen.Dimension.Width/4", "The (initial) x-position of the window")
	flag.StringVar(&c.UI.Window.InitialPositionY, "ui-init-pos-y", "0", "The (initial) y-position of the window")
	flag.Float64Var(&c.UI.Window.InitialZoom, "ui-init-zoom", 1.0, "The (initial) zoom level of the window")
	flag.UintVar(&c.UI.Window.BackgroundColor.R, "ui-bg-color-r", 255, "The window's background color (red value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.G, "ui-bg-color-g", 255, "The window's background color (green value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.B, "ui-bg-color-b", 255, "The window's background color (blue value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.A, "ui-bg-color-a", 192, "The window's background color (alpha value)")
	flag.IntVar(&c.UI.Window.StartState, "ui-start-state", int(options.Normal), fmt.Sprintf("The window start state (normal(%d), minimized(%d), maximized(%d), fullscreen(%d))", options.Normal, options.Minimised, options.Maximised, options.Fullscreen))
	flag.BoolVar(&c.UI.Window.Frameless, "ui-frameless", true, "Should the window be frameless")
	flag.BoolVar(&c.UI.Window.Resizeable, "ui-resizeable", true, "Should the window be resizeable")
	flag.StringVar(&c.UI.Window.Translucent, "ui-translucent", TranslucentHover, fmt.Sprintf("When the window should be translucent (%s, %s, %s)", TranslucentNever, TranslucentEver, TranslucentHover))
	flag.StringVar(&c.UI.QuitShortcut.Code, "ui-quit-shortcut-keycode", "escape", "The shortcut for quitting the application (key-code)")
	flag.BoolVar(&c.UI.QuitShortcut.Ctrl, "ui-quit-shortcut-ctrl", false, "The shortcut for quitting the application (control-key must be pressed)")
	flag.BoolVar(&c.UI.QuitShortcut.Shift, "ui-quit-shortcut-shift", false, "The shortcut for quitting the application (shift-key must be pressed)")
	flag.BoolVar(&c.UI.QuitShortcut.Alt, "ui-quit-shortcut-alt", false, "The shortcut for quitting the application (alt-key must be pressed)")
	flag.BoolVar(&c.UI.QuitShortcut.Meta, "ui-quit-shortcut-meta", false, "The shortcut for quitting the application (meta-key must be pressed)")
	flag.StringVar(&c.UI.Theme, "ui-theme", ThemeSystem, fmt.Sprintf("The theme to use ('%s', '%s', '%s')", ThemeLight, ThemeDark, ThemeSystem))
	flag.StringVar(&c.UI.CodeStyle, "ui-code-style", "github", "The code style to use")
	flag.StringVar(&c.UI.Language, "ui-lang", os.Getenv("LANG"), "The language to use")

	flag.StringVar(&c.Backend, "backend", BackendCopilot, fmt.Sprintf("The backend to use ('%s', '%s', '%s', '%s', '%s', '%s')", BackendCopilot, BackendOpenAI, BackendAnythingLLM, BackendOllama, BackendMistral, BackendAnthropic))

	configureLocalai(&c.LocalAI)
	configureOpenai(&c.OpenAI)
	configureAnythingLLM(&c.AnythingLLM)
	configureOllama(&c.Ollama)
	configureMistral(&c.Mistral)
	configureAnthropic(&c.Anthropic)
	configureCallOptions(&c.CallOptions)

	flag.StringVar(&c.Printer.Format, "print-format", PrinterFormatJSON, fmt.Sprintf("Response printer format (%s, %s)", PrinterFormatPlain, PrinterFormatJSON))
	flag.StringVar(&c.Printer.targets, "print-targets", PrinterTargetOut, fmt.Sprintf("Comma seperated response printer targets (%s, %s, <path/to/file>)", PrinterTargetOut, PrinterTargetErr))

	flag.Usage = func() {
		printUsage(flag.CommandLine.Output())
	}

	flag.CommandLine.Parse(arguments)

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

	if c.UI.Window.BackgroundColor.R > 255 {
		return fmt.Errorf("Invalid background color (red)")
	}
	if c.UI.Window.BackgroundColor.G > 255 {
		return fmt.Errorf("Invalid background color (green)")
	}
	if c.UI.Window.BackgroundColor.B > 255 {
		return fmt.Errorf("Invalid background color (blue)")
	}
	if c.UI.Window.BackgroundColor.A > 255 {
		return fmt.Errorf("Invalid background color (alpha)")
	}
	if c.UI.Theme != ThemeDark && c.UI.Theme != ThemeLight && c.UI.Theme != ThemeSystem {
		return fmt.Errorf("Invalid theme")
	}
	if c.UI.Window.StartState < int(options.Normal) || c.UI.Window.StartState > int(options.Fullscreen) {
		return fmt.Errorf("Invalid window start state")
	}
	if c.UI.Window.InitialZoom < 0.1 || c.UI.Window.InitialZoom > 10.0 {
		return fmt.Errorf("Invalid window zoom: value must be between 0.1 and 10.0")
	}

	if c.UI.Window.MaxHeight != "" {
		if err := ValidateExpression(c.UI.Window.MaxHeight); err != nil {
			return fmt.Errorf("Invalid window max height expression: %w", err)
		}
	}

	if c.UI.Window.InitialWidth != "" {
		if err := ValidateExpression(c.UI.Window.InitialWidth); err != nil {
			return fmt.Errorf("Invalid window initial width expression: %w", err)
		}
	}

	if c.UI.Window.InitialPositionX != "" {
		if err := ValidateExpression(c.UI.Window.InitialPositionX); err != nil {
			return fmt.Errorf("Invalid window initial x-position expression: %w", err)
		}
	}

	if c.UI.Window.InitialPositionY != "" {
		if err := ValidateExpression(c.UI.Window.InitialPositionY); err != nil {
			return fmt.Errorf("Invalid window initial y-position expression: %w", err)
		}
	}

	if c.UI.Window.Translucent != TranslucentNever && c.UI.Window.Translucent != TranslucentEver && c.UI.Window.Translucent != TranslucentHover {
		return fmt.Errorf("Invalid window translucent value")
	}

	switch c.Backend {
	case BackendCopilot:
		if !llms.IsCopilotInstalled() {
			return fmt.Errorf("GitHub Copilot is not installed")
		}
	case BackendLocalAI:
		if ve := c.LocalAI.Validate(); ve != nil {
			return ve
		}
	case BackendOpenAI:
		if ve := c.OpenAI.Validate(); ve != nil {
			return ve
		}
	case BackendAnythingLLM:
		if ve := c.AnythingLLM.Validate(); ve != nil {
			return ve
		}
	case BackendOllama:
		if ve := c.Ollama.Validate(); ve != nil {
			return ve
		}
	case BackendMistral:
		if ve := c.Mistral.Validate(); ve != nil {
			return ve
		}
	case BackendAnthropic:
		if ve := c.Anthropic.Validate(); ve != nil {
			return ve
		}
	default:
		return fmt.Errorf("Invalid backend")
	}

	if ve := c.CallOptions.Validate(); ve != nil {
		return ve
	}

	if c.Printer.Format != PrinterFormatJSON && c.Printer.Format != PrinterFormatPlain {
		return fmt.Errorf("Invalid response printer format")
	}

	return nil
}
