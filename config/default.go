package config

import (
	"github.com/rainu/ask-mai/config/llm"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
)

func defaultConfig() *Config {
	return &Config{
		LogLevel: int(slog.LevelError),
		LLM: llm.LLMConfig{
			Backend: llm.BackendCopilot,
			CallOptions: llm.CallOptionsConfig{
				Temperature: -1,
				TopK:        -1,
				TopP:        -1,
			},
			OpenAI: llm.OpenAIConfig{
				APIType: string(openai.APITypeOpenAI),
				Model:   "gpt-4o-mini",
			},
		},
		UI: UIConfig{
			Prompt: PromptConfig{
				MinRows:         1,
				MaxRows:         4,
				SubmitShortcut:  Shortcut{Alt: true, Code: "enter"},
				PinTop:          true,
				InitAttachments: []string{},
			},
			FileDialog: FileDialogConfig{
				ShowHiddenFiles:            true,
				ResolveAliases:             true,
				TreatPackagesAsDirectories: true,
				FilterDisplay:              []string{},
				FilterPattern:              []string{},
			},
			Window: WindowConfig{
				Title:            "Prompt - Ask mAI",
				InitialWidth:     ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Width/2"},
				MaxHeight:        ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Height/3"},
				InitialPositionX: ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Width/4"},
				InitialPositionY: ExpressionContainer{Expression: "0"},
				InitialZoom:      ExpressionContainer{Expression: "1.0"},
				BackgroundColor:  WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192},
				StartState:       int(options.Normal),
				Translucent:      TranslucentHover,
				Frameless:        true,
				AlwaysOnTop:      true,
				Resizeable:       true,
			},
			QuitShortcut: Shortcut{Code: "escape"},
			Theme:        ThemeSystem,
			CodeStyle:    "github",
			Language:     os.Getenv("LANG"),
		},
		Printer: PrinterConfig{
			Format:     PrinterFormatJSON,
			Targets:    []io.WriteCloser{os.Stdout},
			TargetsRaw: []string{PrinterTargetOut},
		},
	}
}
