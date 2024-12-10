package config

import (
	"context"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
)

type NoopLooger struct {
}

func (n NoopLooger) Print(message string) {
}

func (n NoopLooger) Trace(message string) {
}

func (n NoopLooger) Debug(message string) {
}

func (n NoopLooger) Info(message string) {
}

func (n NoopLooger) Warning(message string) {
}

func (n NoopLooger) Error(message string) {
}

func (n NoopLooger) Fatal(message string) {
}

func printUsage(output io.Writer) {
	fmt.Fprintf(output, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Fprintf(output, "\nAvailable Expression variables:\n")

	wails.Run(&options.App{
		StartHidden: true,
		Frameless:   true,
		Width:       1,
		Height:      1,
		OnStartup: func(ctx context.Context) {
			screens, err := runtime.ScreenGetAll(ctx)
			if err != nil {
				println(err.Error())
				return
			}

			flatMap, err := FromScreens(screens).ToFlatMap()
			if err != nil {
				println(err.Error())
				return
			}
			for k, v := range flatMap {
				fmt.Fprintf(flag.CommandLine.Output(), "  - %s: %v\n", k, v)
			}
		},
		OnDomReady: func(ctx context.Context) {
			runtime.Quit(ctx)
		},
		AssetServer: &assetserver.Options{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		},
		Logger:           &NoopLooger{},
		WindowStartState: options.Minimised,
	})
}
