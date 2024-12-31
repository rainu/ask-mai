package controller

import (
	"embed"
	"fmt"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
	"github.com/rainu/ask-mai/llms"
	langChainLLM "github.com/tmc/langchaingo/llms"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log/slog"
)

func BuildFromConfig(cfg *config.Config) (ctrl *Controller, err error) {
	ctrl = &Controller{
		appConfig: cfg,
	}

	printer := io.MultiResponsePrinter{}
	if cfg.Printer.Format == config.PrinterFormatPlain {
		for _, target := range cfg.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.PlainResponsePrinter{Target: target})
		}
	} else if cfg.Printer.Format == config.PrinterFormatJSON {
		for _, target := range cfg.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.JsonResponsePrinter{Target: target})
		}
	}
	ctrl.printer = printer

	switch cfg.Backend {
	case config.BackendCopilot:
		ctrl.aiModel, err = llms.NewCopilot()
	case config.BackendLocalAI:
		ctrl.aiModel, err = llms.NewLocalAI(
			cfg.LocalAI.AsOptions(),
		)
	case config.BackendOpenAI:
		ctrl.aiModel, err = llms.NewOpenAI(
			cfg.OpenAI.AsOptions(),
		)
	case config.BackendAnythingLLM:
		ctrl.aiModel, err = llms.NewAnythingLLM(
			cfg.AnythingLLM.BaseURL,
			cfg.AnythingLLM.Token,
			cfg.AnythingLLM.Workspace,
		)
	case config.BackendOllama:
		ctrl.aiModel, err = llms.NewOllama(
			cfg.Ollama.AsOptions(),
		)
	case config.BackendMistral:
		ctrl.aiModel, err = llms.NewMistral(
			cfg.Mistral.AsOptions(),
		)
	case config.BackendAnthropic:
		ctrl.aiModel, err = llms.NewAnthropic(
			cfg.Anthropic.AsOptions(),
		)
	default:
		err = fmt.Errorf("unknown backend: %s", cfg.Backend)
	}

	if err != nil {
		err = fmt.Errorf("error creating ai model: %w", err)
		return
	}

	if cfg.UI.Prompt.InitValue != "" {
		// ask the model the first question in background
		go ctrl.LLMAsk(LLMAskArgs{
			History: LLMMessages{{
				Content: cfg.UI.Prompt.InitValue,
				Role:    string(langChainLLM.ChatMessageTypeHuman),
			}},
		})
	}

	return
}

func GetOptions(c *Controller, icon []byte, assets embed.FS) *options.App {
	ac := c.appConfig
	translucent := true
	if ac.UI.Window.Translucent == config.TranslucentNever {
		translucent = false
	}

	return &options.App{
		Title:             ac.UI.Window.Title,
		Height:            1,
		DisableResize:     !ac.UI.Window.Resizeable,
		Fullscreen:        false,
		Frameless:         ac.UI.Window.Frameless,
		StartHidden:       true,
		HideWindowOnClose: false,
		AlwaysOnTop:       ac.UI.Window.AlwaysOnTop,
		BackgroundColour: &options.RGBA{
			R: uint8(ac.UI.Window.BackgroundColor.R),
			G: uint8(ac.UI.Window.BackgroundColor.G),
			B: uint8(ac.UI.Window.BackgroundColor.B),
			A: uint8(ac.UI.Window.BackgroundColor.A),
		},
		Menu:   nil,
		Logger: newDefaultLogger(),
		LogLevel: func() logger.LogLevel {
			switch ac.LogLevel {
			case int(slog.LevelDebug):
				return logger.DEBUG
			case int(slog.LevelInfo):
				return logger.INFO
			case int(slog.LevelWarn):
				return logger.WARNING
			case int(slog.LevelError):
				return logger.ERROR
			default:
				return logger.ERROR
			}
		}(),
		OnStartup:        c.startup,
		OnDomReady:       c.domReady,
		OnBeforeClose:    c.beforeClose,
		OnShutdown:       c.shutdown,
		WindowStartState: options.WindowStartState(ac.UI.Window.StartState),
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    nil,
			Middleware: nil,
		},
		Bind: []interface{}{
			c,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               translucent,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme: func() windows.Theme {
				switch ac.UI.Theme {
				case config.ThemeDark:
					return windows.Dark
				case config.ThemeLight:
					return windows.Light
				default:
					return windows.SystemDefault
				}
			}(),
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  translucent,
			About: &mac.AboutInfo{
				Title:   ac.UI.Window.Title,
				Message: "Ask mAI is a simple application to ask questions to AI models.",
				Icon:    icon,
			},
		},
		Linux: &linux.Options{
			Icon:                icon,
			WindowIsTranslucent: translucent,
			ProgramName:         ac.UI.Window.Title,
		},
	}
}
