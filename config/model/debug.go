package model

import (
	"fmt"
	"log/slog"
)

type DebugConfig struct {
	LogLevel              int                   `yaml:"log-level,omitempty"`
	PprofAddress          string                `yaml:"pprof-address,omitempty" usage:"Address for the pprof server (only available for debug binary)"`
	VueDevTools           VueDevToolsConfig     `yaml:"vue-dev-tools,omitempty"`
	WebKit                WebKitInspectorConfig `yaml:"webkit,omitempty" usage:"Webkit debug configuration (only available for debug binary): "`
	DisableCrashDetection bool                  `yaml:"disable-crash-detection,omitempty" usage:"Disable crash detection. If a crash is detected the application will try to recover the last state"`
}

type VueDevToolsConfig struct {
	Host string `yaml:"host,omitempty" usage:"The host of the vue dev tools server. If empty the dev tools will be disabled"`
	Port int    `yaml:"port,omitempty" usage:"The port of the vue dev tools server"`
}

type WebKitInspectorConfig struct {
	OpenInspectorOnStartup bool   `yaml:"open-inspector,omitempty" usage:"Open the inspector on startup"`
	HttpServerAddress      string `yaml:"http-server,omitempty" usage:"Starts a http server for the inspector"`
}

func (d *DebugConfig) SetDefaults() {
	d.LogLevel = int(slog.LevelError)
	d.PprofAddress = ":6060"
}

func (v *VueDevToolsConfig) SetDefaults() {
	v.Port = 8098
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

	return nil
}
