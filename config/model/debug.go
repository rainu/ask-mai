package model

import (
	"fmt"
	"github.com/rainu/go-yacl"
	"log/slog"
	"strings"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type DebugConfig struct {
	LogLevel              string                `yaml:"log-level,omitempty"`
	LogLevelParsed        *slog.Level           `yaml:"-"`
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
	if d.LogLevel == "" {
		d.LogLevel = "error"
		d.LogLevelParsed = yacl.P(slog.LevelError)
	}
	if d.PprofAddress == "" {
		d.PprofAddress = ":6060"
	}
}

func (v *VueDevToolsConfig) SetDefaults() {
	if v.Port == 0 {
		v.Port = 8098
	}
}

func (d *DebugConfig) GetUsage(field string) string {
	switch field {
	case "LogLevel":
		return fmt.Sprintf("Log level (%s, %s, %s, %s)", LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
	}
	return ""
}

func (d *DebugConfig) Validate() error {
	switch strings.ToLower(d.LogLevel) {
	case LogLevelDebug:
		d.LogLevelParsed = yacl.P(slog.LevelDebug)
	case LogLevelInfo:
		d.LogLevelParsed = yacl.P(slog.LevelInfo)
	case LogLevelWarn:
		d.LogLevelParsed = yacl.P(slog.LevelWarn)
	case LogLevelError:
		d.LogLevelParsed = yacl.P(slog.LevelError)
	default:
		return fmt.Errorf("Invalid log level '%s'", d.LogLevel)
	}

	return nil
}
