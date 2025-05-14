package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/config/model"
	"github.com/rainu/ask-mai/controller"
	"github.com/rainu/ask-mai/health"
	"github.com/rainu/ask-mai/notification"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/wailsapp/wails/v2"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"strings"
)

const (
	lastStateEnv = "_ASK_MAI_LAST_STATE"
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
	notification.SetNotificationIcon(icon)
}

// this function will be set by debug.go:init() - if "debug" flag is available
var onStartUp func(c *model.Config)

func main() {
	buildMode := slices.ContainsFunc(os.Environ(), func(s string) bool {
		return strings.HasPrefix(s, "tsprefix=")
	})

	cfg := config.Parse(os.Args[1:], os.Environ())
	if cfg.Version {
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

		slog.SetLogLoggerLevel(*cfg.DebugConfig.LogLevelParsed)
		if onStartUp != nil {
			onStartUp(cfg)
		}
	}
	defer cfg.MainProfile.Printer.Close()

	ctrl, err := controller.BuildFromConfig(cfg, os.Getenv(lastStateEnv), buildMode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
		return
	}

	// Create application with options
	if cfg.DebugConfig.WebKit.HttpServerAddress != "" {
		// the underlying webview library will use this environment variable to start the inspector server
		os.Setenv("WEBKIT_INSPECTOR_HTTP_SERVER", cfg.DebugConfig.WebKit.HttpServerAddress)
	}

	if !cfg.DebugConfig.DisableCrashDetection {
		oCtx, oCancel := context.WithCancel(context.Background())
		health.ObserveProcess(oCtx, 98.0, func() {
			if ctrl.IsAppMounted() {
				slog.Warn("Restarting application because of high CPU usage: Seems like a freeze.")
				notification.Notify("ask-mai", "Restarting application!")

				ctrl.TriggerRestart()
				oCancel() //prevent multiple restarts
			}
		})
	}

	err = wails.Run(controller.GetOptions(ctrl, icon, assets))
	if !buildMode && ctrl.GetLastState() != "" {
		ae := map[any]any{
			lastStateEnv: ctrl.GetLastState(),
		}

		cmdchain.Builder().WithInput(os.Stdin).
			Join(os.Args[0], os.Args[1:]...).WithAdditionalEnvironmentMap(ae).
			Finalize().WithError(os.Stderr).WithOutput(os.Stdout).
			Run()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
		return
	}
}
