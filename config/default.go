package config

import (
	"github.com/kirsle/configdir"
	"github.com/rainu/ask-mai/config/model"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/config/model/llm"
	"github.com/rainu/ask-mai/config/model/llm/tools"
	"github.com/rainu/ask-mai/expression"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime"
)

func defaultConfig() (result *model.Config) {
	confPath := configdir.LocalConfig("ask-mai")
	confPath = path.Join(confPath, "history")

	result = &model.Config{
		Debug: model.DebugConfig{
			LogLevel:     int(slog.LevelError),
			PprofAddress: ":6060",
			VueDevTools: model.VueDevToolsConfig{
				Host: "",
				Port: 8098,
			},
			WebKit: model.WebKitInspectorConfig{
				OpenInspectorOnStartup: false,
				HttpServerAddress:      "",
			},
			RestartShortcut: model.Shortcut{Binding: []string{"Alt+KeyR"}},
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
				Token: common.Secret{Command: common.SecretCommand{Args: []string{}}},
				Model: "claude-3-5-haiku-latest",
			},
			AnythingLLM: llm.AnythingLLMConfig{
				Token: common.Secret{Command: common.SecretCommand{Args: []string{}}},
				Thread: llm.AnythingLLMThreadConfig{
					Delete: false,
					Name: common.StringContainer{
						Expression: `'ask-mai - ' + new Date().toISOString()`,
					},
				},
			},
			DeepSeek: llm.DeepSeekConfig{
				APIKey: common.Secret{Command: common.SecretCommand{Args: []string{}}},
				Model:  "deepseek-chat",
			},
			Mistral: llm.MistralConfig{
				ApiKey: common.Secret{Command: common.SecretCommand{Args: []string{}}},
			},
			LocalAI: llm.LocalAIConfig{
				APIKey: common.Secret{Command: common.SecretCommand{Args: []string{}}},
			},
			OpenAI: llm.OpenAIConfig{
				APIKey:  common.Secret{Command: common.SecretCommand{Args: []string{}}},
				APIType: string(openai.APITypeOpenAI),
				Model:   "gpt-4o-mini",
			},
		},
		UI: model.UIConfig{
			Prompt: model.PromptConfig{
				MinRows:         1,
				MaxRows:         4,
				SubmitShortcut:  model.Shortcut{Binding: []string{"Alt+Enter", "Alt+NumpadEnter"}},
				PinTop:          true,
				InitAttachments: []string{},
			},
			FileDialog: model.FileDialogConfig{
				ShowHiddenFiles:            true,
				ResolveAliases:             true,
				TreatPackagesAsDirectories: true,
				FilterDisplay:              []string{},
				FilterPattern:              []string{},
			},
			Window: model.WindowConfig{
				Title:            "Prompt - Ask mAI",
				InitialWidth:     common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/2"},
				MaxHeight:        common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Height/3"},
				InitialPositionX: common.NumberContainer{Expression: expression.VarNameScreens + ".CurrentScreen.Dimension.Width/4"},
				InitialPositionY: common.NumberContainer{Expression: "0"},
				InitialZoom:      common.NumberContainer{Expression: "1.0"},
				BackgroundColor:  model.WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192},
				StartState:       int(options.Normal),
				Translucent:      model.TranslucentHover,
				Frameless:        true,
				AlwaysOnTop:      true,
				Resizeable:       true,
			},
			QuitShortcut: model.Shortcut{Binding: []string{"Escape"}},
			Theme:        model.ThemeSystem,
			CodeStyle:    "github",
			Language:     os.Getenv("LANG"),
		},
		History: model.History{
			Path: confPath,
		},
		Printer: model.PrinterConfig{
			Format:     model.PrinterFormatJSON,
			Targets:    []io.WriteCloser{os.Stdout},
			TargetsRaw: []string{model.PrinterTargetOut},
		},
	}

	if runtime.GOOS == "windows" {
		result.UI.Window.TitleBarHeight = 32
	} else if runtime.GOOS == "darwin" {
		result.UI.Window.TitleBarHeight = 28
	}

	if llms.IsCopilotInstalled() {
		result.LLM.Backend = "copilot"
	}

	return
}
