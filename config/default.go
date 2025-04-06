package config

import (
	"github.com/kirsle/configdir"
	"github.com/rainu/ask-mai/config/common"
	"github.com/rainu/ask-mai/config/llm"
	"github.com/rainu/ask-mai/config/llm/tools"
	"github.com/rainu/ask-mai/expression"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
	"path"
)

func defaultConfig() (result *Config) {
	confPath := configdir.LocalConfig("ask-mai")
	confPath = path.Join(confPath, "history")

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
					Name: common.StringContainer{
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
				InitialWidth:     common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/2"},
				MaxHeight:        common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Height/3"},
				InitialPositionX: common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/4"},
				InitialPositionY: common.NumberContainer{Expression: "0"},
				InitialZoom:      common.NumberContainer{Expression: "1.0"},
				BackgroundColor:  WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192},
				StartState:       int(options.Normal),
				Translucent:      TranslucentHover,
				Frameless:        true,
				AlwaysOnTop:      true,
				Resizeable:       true,
			},
			QuitShortcut: Shortcut{Binding: []string{"Escape"}},
			Theme:        ThemeSystem,
			CodeStyle:    "github",
			Language:     os.Getenv("LANG"),
		},
		History: History{
			Path: confPath,
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
