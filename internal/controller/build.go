package controller

import (
	"embed"
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/io"
	langChainLLM "github.com/tmc/langchaingo/llms"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log/slog"
	"net/http"
	"time"
)

func BuildFromConfig(cfg *model.Config, lastState string, buildMode bool) (ctrl *Controller, err error) {
	activeCfg := cfg.GetActiveProfile()
	ctrl = &Controller{
		appConfig: cfg,
		lastState: lastState,
	}

	if buildMode {
		// in build mode the bindings will be generated
		return ctrl, nil
	}

	printer := io.MultiResponsePrinter{}
	if activeCfg.Printer.Format == model.PrinterFormatPlain {
		for _, target := range activeCfg.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.PlainResponsePrinter{Target: target})
		}
	} else if activeCfg.Printer.Format == model.PrinterFormatJSON {
		for _, target := range activeCfg.Printer.Targets {
			printer.Printers = append(printer.Printers, &io.JsonResponsePrinter{Target: target})
		}
	}
	ctrl.printer = printer

	ctrl.aiModel, err = activeCfg.LLM.BuildLLM()
	if err != nil {
		err = fmt.Errorf("error creating ai model: %w", err)
		return
	}

	if activeCfg.LLM.CallOptions.Prompt.InitValue != "" && lastState == "" {
		args := buildInitialChat(activeCfg)
		go ctrl.LLMAsk(args)
	}

	return
}

func buildInitialChat(activeCfg *model.Profile) LLMAskArgs {
	now := time.Now().Unix()
	// ask the model the first question in background

	args := LLMAskArgs{}

	if activeCfg.LLM.CallOptions.Prompt.System != "" {
		// add the system prompt to the history
		args.History = append(args.History, LLMMessage{
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: activeCfg.LLM.CallOptions.Prompt.System,
			}},
			Role:    string(langChainLLM.ChatMessageTypeSystem),
			Created: now,
		})
	}

	for _, message := range activeCfg.LLM.CallOptions.Prompt.InitMessages {
		args.History = append(args.History, LLMMessage{
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: message.Content,
			}},
			Role:    string(message.Role),
			Created: now,
		})
	}

	if activeCfg.LLM.CallOptions.Prompt.InitValue != "" {
		message := LLMMessage{
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: activeCfg.LLM.CallOptions.Prompt.InitValue,
			}},
			Role:    string(langChainLLM.ChatMessageTypeHuman),
			Created: now,
		}

		for _, attachment := range activeCfg.LLM.CallOptions.Prompt.InitAttachments {
			message.ContentParts = append(message.ContentParts, LLMMessageContentPart{
				Type:    LLMMessageContentPartTypeAttachment,
				Content: attachment,
			})
		}
		args.History = append(args.History, message)
	}

	return args
}

func GetOptions(c *Controller, icon []byte, assets embed.FS) *options.App {
	ac := c.getProfile()
	translucent := true
	if ac.UI.Window.Translucent == model.TranslucentNever {
		translucent = false
	}

	return &options.App{
		Title:             ac.UI.Window.Title,
		Height:            1,
		DisableResize:     !*ac.UI.Window.Resizeable,
		Fullscreen:        false,
		Frameless:         *ac.UI.Window.Frameless,
		StartHidden:       true,
		HideWindowOnClose: false,
		AlwaysOnTop:       *ac.UI.Window.AlwaysOnTop,
		BackgroundColour: &options.RGBA{
			R: uint8(*ac.UI.Window.BackgroundColor.R),
			G: uint8(*ac.UI.Window.BackgroundColor.G),
			B: uint8(*ac.UI.Window.BackgroundColor.B),
			A: uint8(*ac.UI.Window.BackgroundColor.A),
		},
		Menu: nil,
		Debug: options.Debug{
			OpenInspectorOnStartup: c.appConfig.DebugConfig.WebKit.OpenInspectorOnStartup,
		},
		Logger: newDefaultLogger(),
		LogLevel: func() logger.LogLevel {
			switch *c.appConfig.DebugConfig.LogLevelParsed {
			case slog.LevelDebug:
				return logger.TRACE
			case slog.LevelInfo:
				return logger.INFO
			case slog.LevelWarn:
				return logger.WARNING
			case slog.LevelError:
				return logger.ERROR
			default:
				return logger.ERROR
			}
		}(),
		OnStartup:        c.startup,
		OnDomReady:       c.domReady,
		OnBeforeClose:    c.beforeClose,
		OnShutdown:       c.shutdown,
		WindowStartState: options.WindowStartState(*ac.UI.Window.StartState),
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    http.HandlerFunc(c.handleAsset),
			Middleware: nil,
		},
		Bind: []interface{}{
			c,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               translucent,
			DisableWindowIcon:                 !*ac.UI.Window.ShowTitleBar,
			DisableFramelessWindowDecorations: !*ac.UI.Window.ShowTitleBar,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme: func() windows.Theme {
				switch ac.UI.Theme {
				case model.ThemeDark:
					return windows.Dark
				case model.ThemeLight:
					return windows.Light
				default:
					return windows.SystemDefault
				}
			}(),
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  !*ac.UI.Window.ShowTitleBar,
				HideTitleBar:               !*ac.UI.Window.ShowTitleBar,
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
