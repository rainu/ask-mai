package controller

import (
	"embed"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/backend/anythingllm"
	"github.com/rainu/ask-mai/backend/copilot"
	"github.com/rainu/ask-mai/backend/openai"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log/slog"
)

func BuildFromConfig(cfg *config.Config) (ctrl *Controller) {
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
		ctrl.backendBuilder = backend.Builder{
			Build: copilot.NewCopilot,
			Type:  backend.TypeSingleShot,
		}
	case config.BackendOpenAI:
		ctrl.backendBuilder = backend.Builder{
			Build: func() (backend.Handle, error) {
				return openai.NewOpenAI(cfg.OpenAI.APIKey, cfg.OpenAI.SystemPrompt)
			},
			Type: backend.TypeMultiShot,
		}
	case config.BackendAnythingLLM:
		ctrl.backendBuilder = backend.Builder{
			Build: func() (backend.Handle, error) {
				return anythingllm.NewAnythingLLM(cfg.AnythingLLM.BaseURL, cfg.AnythingLLM.Token, cfg.AnythingLLM.Workspace)
			},
			Type: backend.TypeMultiShot,
		}
	}

	return
}

func GetOptions(c *Controller, icon []byte, assets embed.FS) *options.App {
	ac := c.appConfig
	return &options.App{
		Title:             ac.UI.Window.Title,
		Height:            1,
		DisableResize:     !ac.UI.Window.Resizeable,
		Fullscreen:        false,
		Frameless:         ac.UI.Window.Frameless,
		StartHidden:       true,
		HideWindowOnClose: false,
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
			WindowIsTranslucent:               false,
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
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   ac.UI.Window.Title,
				Message: "Ask mAI is a simple application to ask questions to AI models.",
				Icon:    icon,
			},
		},
		Linux: &linux.Options{
			Icon: icon,
		},
	}
}
