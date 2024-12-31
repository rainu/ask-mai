package config

import (
	"errors"
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

	LogLevel     int
	PrintVersion bool
}

type UIConfig struct {
	Window       WindowConfig
	Prompt       PromptConfig
	Stream       bool
	QuitShortcut Shortcut
	Theme        string
	CodeStyle    string
	Language     string
}

type PromptConfig struct {
	InitValue      string
	MinRows        uint
	MaxRows        uint
	SubmitShortcut Shortcut
}

type WindowConfig struct {
	Title            string
	InitialWidth     ExpressionContainer
	MaxHeight        ExpressionContainer
	InitialPositionX ExpressionContainer
	InitialPositionY ExpressionContainer
	InitialZoom      ExpressionContainer
	BackgroundColor  WindowBackgroundColor
	StartState       int
	AlwaysOnTop      bool
	Frameless        bool
	Resizeable       bool
	Translucent      string
}

type ExpressionContainer struct {
	Expression string
	Value      float64
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

	flag.StringVar(&c.UI.Prompt.InitValue, "ui-prompt-value", "", "The (initial) prompt to use")
	flag.UintVar(&c.UI.Prompt.MinRows, "ui-prompt-min-rows", 1, "The minimal number of rows the prompt should have")
	flag.UintVar(&c.UI.Prompt.MaxRows, "ui-prompt-max-rows", 4, "The maximal number of rows the prompt should have")
	flag.StringVar(&c.UI.Prompt.SubmitShortcut.Code, "ui-prompt-submit-key", "enter", "The shortcut for submit the prompt (key-code)")
	flag.BoolVar(&c.UI.Prompt.SubmitShortcut.Alt, "ui-prompt-submit-alt", true, "The shortcut for submit the prompt (alt-key must be pressed)")
	flag.BoolVar(&c.UI.Prompt.SubmitShortcut.Ctrl, "ui-prompt-submit-ctrl", false, "The shortcut for submit the prompt (control-key must be pressed)")
	flag.BoolVar(&c.UI.Prompt.SubmitShortcut.Meta, "ui-prompt-submit-meta", false, "The shortcut for submit the prompt (meta-key must be pressed)")
	flag.BoolVar(&c.UI.Prompt.SubmitShortcut.Shift, "ui-prompt-submit-shift", false, "The shortcut for submit the prompt (shift-key must be pressed)")
	flag.BoolVar(&c.UI.Stream, "ui-stream", false, "Should the output be streamed")
	flag.StringVar(&c.UI.Window.Title, "ui-title", "Prompt - Ask mAI", "The window title")
	flag.StringVar(&c.UI.Window.InitialWidth.Expression, "ui-init-width", "v.CurrentScreen.Dimension.Width/2", "Expression: The (initial) width of the window")
	flag.StringVar(&c.UI.Window.MaxHeight.Expression, "ui-max-height", "v.CurrentScreen.Dimension.Height/3", "Expression: The maximal height of the chat response area")
	flag.StringVar(&c.UI.Window.InitialPositionX.Expression, "ui-init-pos-x", "v.CurrentScreen.Dimension.Width/4", "Expression: The (initial) x-position of the window")
	flag.StringVar(&c.UI.Window.InitialPositionY.Expression, "ui-init-pos-y", "0", "Expression: The (initial) y-position of the window")
	flag.StringVar(&c.UI.Window.InitialZoom.Expression, "ui-init-zoom", "1.0", "Expression: The (initial) zoom level of the window")
	flag.UintVar(&c.UI.Window.BackgroundColor.R, "ui-bg-color-r", 255, "The window's background color (red value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.G, "ui-bg-color-g", 255, "The window's background color (green value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.B, "ui-bg-color-b", 255, "The window's background color (blue value)")
	flag.UintVar(&c.UI.Window.BackgroundColor.A, "ui-bg-color-a", 192, "The window's background color (alpha value)")
	flag.IntVar(&c.UI.Window.StartState, "ui-start-state", int(options.Normal), fmt.Sprintf("The window start state (normal(%d), minimized(%d), maximized(%d), fullscreen(%d))", options.Normal, options.Minimised, options.Maximised, options.Fullscreen))
	flag.BoolVar(&c.UI.Window.Frameless, "ui-frameless", true, "Should the window be frameless")
	flag.BoolVar(&c.UI.Window.AlwaysOnTop, "ui-always-on-top", true, "Should the window be always on top")
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

	flag.BoolVar(&c.PrintVersion, "v", false, "Show the version")

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

func (c *Config) Validate() error {
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

	if c.UI.Window.MaxHeight.Expression != "" {
		if err := ValidateExpression(c.UI.Window.MaxHeight.Expression); err != nil {
			return fmt.Errorf("Invalid window max height expression: %w", err)
		}
	}

	if c.UI.Window.InitialWidth.Expression != "" {
		if err := ValidateExpression(c.UI.Window.InitialWidth.Expression); err != nil {
			return fmt.Errorf("Invalid window initial width expression: %w", err)
		}
	}

	if c.UI.Window.InitialPositionX.Expression != "" {
		if err := ValidateExpression(c.UI.Window.InitialPositionX.Expression); err != nil {
			return fmt.Errorf("Invalid window initial x-position expression: %w", err)
		}
	}

	if c.UI.Window.InitialPositionY.Expression != "" {
		if err := ValidateExpression(c.UI.Window.InitialPositionY.Expression); err != nil {
			return fmt.Errorf("Invalid window initial y-position expression: %w", err)
		}
	}

	if c.UI.Window.InitialZoom.Expression != "" {
		if err := ValidateExpression(c.UI.Window.InitialZoom.Expression); err != nil {
			return fmt.Errorf("Invalid window initial zoom expression: %w", err)
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

func (c *Config) ResolveExpressions(variables Variables) (err error) {
	var curErr error

	c.UI.Window.InitialWidth.Value, curErr = Expression(c.UI.Window.InitialWidth.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial width expression: %w", curErr))
	}
	c.UI.Window.MaxHeight.Value, curErr = Expression(c.UI.Window.MaxHeight.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving max width expression: %w", curErr))
	}
	c.UI.Window.InitialPositionX.Value, curErr = Expression(c.UI.Window.InitialPositionX.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial x-position expression: %w", curErr))
	}
	c.UI.Window.InitialPositionY.Value, curErr = Expression(c.UI.Window.InitialPositionY.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial y-position expression: %w", curErr))
	}
	c.UI.Window.InitialZoom.Value, curErr = Expression(c.UI.Window.InitialZoom.Expression).Calculate(variables)
	if curErr != nil {
		err = errors.Join(err, fmt.Errorf("error resolving initial zoom expression: %w", curErr))
	}

	return
}
