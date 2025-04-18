package model

import (
	"fmt"
	"log/slog"
)

type DebugConfig struct {
	LogLevel              int                   `yaml:"log-level"`
	PprofAddress          string                `yaml:"pprof-address" usage:"Address for the pprof server (only available for debug binary)"`
	VueDevTools           VueDevToolsConfig     `yaml:"vue-dev-tools"`
	WebKit                WebKitInspectorConfig `yaml:"webkit" usage:"Webkit debug configuration (only available for debug binary): "`
	DisableCrashDetection bool                  `yaml:"disable-crash-detection" usage:"Disable crash detection. If a crash is detected the application will try to recover the last state"`
	RestartShortcut       Shortcut              `yaml:"restart-shortcut" usage:"The shortcut for triggering a restart: "`
	PrintVersion          bool                  `config:"version" yaml:"-" short:"v" usage:"Show the version"`
}

type VueDevToolsConfig struct {
	Host string `yaml:"host" usage:"The host of the vue dev tools server. If empty the dev tools will be disabled"`
	Port int    `yaml:"port" usage:"The port of the vue dev tools server"`
}

type WebKitInspectorConfig struct {
	OpenInspectorOnStartup bool   `yaml:"open-inspector" usage:"Open the inspector on startup"`
	HttpServerAddress      string `yaml:"http-server" usage:"Starts a http server for the inspector"`
}

func (d *DebugConfig) GetUsage(field string) string {
	switch field {
	case "LogLevel":
		return fmt.Sprintf("Log level (debug(%d), info(%d), warn(%d), error(%d))", slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError)
	}
	return ""
}

func (d *DebugConfig) Validate() error {
	if d.LogLevel < int(slog.LevelDebug) || d.LogLevel > int(slog.LevelError) {
		return fmt.Errorf("Invalid log level")
	}
	if ve := d.RestartShortcut.Validate(); ve != nil {
		return ve
	}

	return nil
}
