package controller

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

type CreateOptionsArg struct {
	Controller *Controller
	Icon       []byte
	Assets     embed.FS

	Width     int
	MaxHeight int
}

func CreateOptions(args CreateOptionsArg) *options.App {
	return &options.App{
		Title:             "Prompt - Ask mAI",
		Width:             args.Width,
		Height:            1,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         true,
		StartHidden:       true,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         args.Controller.startup,
		OnDomReady:        args.Controller.domReady,
		OnBeforeClose:     args.Controller.beforeClose,
		OnShutdown:        args.Controller.shutdown,
		WindowStartState:  options.Normal,
		AssetServer: &assetserver.Options{
			Assets:     args.Assets,
			Handler:    nil,
			Middleware: nil,
		},
		Bind: []interface{}{
			args.Controller,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
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
				Title:   "Ask mAI",
				Message: "Ask mAI is a simple application to ask questions to AI models.",
				Icon:    args.Icon,
			},
		},
		Linux: &linux.Options{
			Icon: args.Icon,
		},
	}
}
