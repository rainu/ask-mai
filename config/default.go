package config

import (
	"github.com/rainu/ask-mai/config/expression"
	"github.com/rainu/ask-mai/config/llm"
	"github.com/rainu/ask-mai/config/llm/tools"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
)

func defaultConfig() (result *Config) {
	result = &Config{
		Debug: DebugConfig{
			LogLevel:     int(slog.LevelError),
			PprofAddress: ":6060",
			VueDevTools: VueDevToolsConfig{
				Host: "",
				Port: 8098,
			},
			WebKit: WebKitInspectorConfig{
				OpenInspectorOnStartup: false,
				HttpServerAddress:      "",
			},
			RestartShortcut: Shortcut{Binding: []string{"Alt+KeyR"}},
			PrintVersion:    false,
		},
		LLM: llm.LLMConfig{
			CallOptions: llm.CallOptionsConfig{
				Temperature: -1,
				TopK:        -1,
				TopP:        -1,
			},
			Tools: tools.Config{
				RawTools: []string{},
				BuiltInTools: tools.BuiltIns{
					CommandExec: tools.CommandExecution{
						NoApprovalCommands:     []string{},
						NoApprovalCommandsExpr: []string{},
					},
				},
			},
			Anthropic: llm.AnthropicConfig{
				Model: "claude-3-5-haiku-latest",
			},
			AnythingLLM: llm.AnythingLLMConfig{
				Thread: llm.AnythingLLMThreadConfig{
					Delete: false,
					Name: expression.StringContainer{
						Expression: `'ask-mai - ' + new Date().toISOString()`,
					},
				},
			},
			DeepSeek: llm.DeepSeekConfig{
				Model: "deepseek-chat",
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
				SubmitShortcut:  Shortcut{Binding: []string{"Alt+Enter", "Alt+NumpadEnter"}},
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
				InitialWidth:     expression.NumberContainer{Expression: "v.CurrentScreen.Dimension.Width/2"},
				MaxHeight:        expression.NumberContainer{Expression: "v.CurrentScreen.Dimension.Height/3"},
				InitialPositionX: expression.NumberContainer{Expression: "v.CurrentScreen.Dimension.Width/4"},
				InitialPositionY: expression.NumberContainer{Expression: "0"},
				InitialZoom:      expression.NumberContainer{Expression: "1.0"},
				BackgroundColor:  WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192},
				StartState:       int(options.Normal),
				Translucent:      TranslucentHover,
				Frameless:        true,
				AlwaysOnTop:      true,
				Resizeable:       true,
			},
			QuitShortcut:   Shortcut{Binding: []string{"Escape"}},
			Theme:          ThemeSystem,
			MinMaxPosition: MinMaxPositionNone,
			CodeStyle:      "github",
			Language:       os.Getenv("LANG"),
		},
		Printer: PrinterConfig{
			Format:     PrinterFormatJSON,
			Targets:    []io.WriteCloser{os.Stdout},
			TargetsRaw: []string{PrinterTargetOut},
		},
	}
	if llms.IsCopilotInstalled() {
		result.LLM.Backend = "copilot"
	}

	return
}
