package main

import (
	"embed"
	"fmt"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/controller"
	"github.com/wailsapp/wails/v2"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"strings"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

// this variable will be set correctly in built-time
var windowMode = "true"

func init() {
	if runtime.GOOS == "windows" && windowMode == "true" {
		// in windows there is no stdout and stderr in window-mode
		// so we need to redirect the log output to a file
		os.Stdout, _ = os.OpenFile("ask-mai.out.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		os.Stderr, _ = os.OpenFile("ask-mai.err.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	}
}

func main() {
	buildMode := slices.ContainsFunc(os.Environ(), func(s string) bool {
		return strings.HasPrefix(s, "tsprefix=")
	})

	cfg := config.Parse(os.Args[1:])
	if cfg.PrintVersion {
		fmt.Fprintln(os.Stderr, versionLine())
		os.Exit(0)
		return
	}
	if !buildMode {
		if err := cfg.Validate(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
			return
		}
	}

	slog.SetLogLoggerLevel(slog.Level(cfg.LogLevel))

	ctrl, err := controller.BuildFromConfig(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
		return
	}

	// Create application with options
	err = wails.Run(controller.GetOptions(ctrl, icon, assets))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
		return
	}
}
