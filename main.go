package main

import (
	"embed"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/controller"
	"github.com/wailsapp/wails/v2"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	buildMode := slices.ContainsFunc(os.Environ(), func(s string) bool {
		return strings.HasPrefix(s, "tsprefix=")
	})

	cfg := config.Parse(os.Args[1:])
	if !buildMode {
		if err := cfg.Validate(); err != nil {
			println(err.Error())
			os.Exit(1)
			return
		}
	}

	slog.SetLogLoggerLevel(slog.Level(cfg.LogLevel))

	ctrl, err := controller.BuildFromConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create application with options
	err = wails.Run(controller.GetOptions(ctrl, icon, assets))

	if err != nil {
		log.Fatal(err)
	}
}
