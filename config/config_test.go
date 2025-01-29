package config

import (
	"github.com/rainu/ask-mai/config/expression"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func modifiedConfig(mod func(*Config)) Config {
	c := defaultConfig()
	mod(c)
	return *c
}

func TestConfig_Parse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      []string
		expected Config
	}{
		{
			name:     "Default values",
			args:     []string{},
			expected: modifiedConfig(func(c *Config) {}),
		},
		{
			name: "Set log level",
			args: []string{"--log-level", "-4"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.LogLevel = int(slog.LevelDebug)
			}),
		},
		{
			name: "Set pprof address",
			args: []string{"--pprof-address", "localhost:6060"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.PprofAddress = "localhost:6060"
			}),
		},
		{
			name: "Set vue dev tools host",
			args: []string{"--vue-dev-tools-host", "localhost"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.VueDevTools.Host = "localhost"
			}),
		},
		{
			name: "Set vue dev tools port",
			args: []string{"--vue-dev-tools-port", "1312"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.VueDevTools.Port = 1312
			}),
		},
		{
			name: "Set open inspector on startup",
			args: []string{"--webkit-open-inspector", "true"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.WebKit.OpenInspectorOnStartup = true
			}),
		},
		{
			name: "Set webkit inspector http server address",
			args: []string{"--webkit-http-server", "127.0.0.1:5000"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.WebKit.HttpServerAddress = "127.0.0.1:5000"
			}),
		},
		{
			name: "Set UI prompt value",
			args: []string{"--ui-prompt-value", "test"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set UI prompt value - shorthand",
			args: []string{"-p", "test"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set system prompt value",
			args: []string{"--call-system-prompt", "test"},
			expected: modifiedConfig(func(c *Config) {
				c.LLM.CallOptions.SystemPrompt = "test"
			}),
		},
		{
			name: "Set UI prompt initial attachments",
			args: []string{"--ui-prompt-attachments", "file1.txt,file2.txt"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set UI prompt initial attachments - shorthand",
			args: []string{"-a", "file1.txt,file2.txt"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set UI prompt min rows",
			args: []string{"--ui-prompt-min-rows", "2"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.MinRows = 2
			}),
		},
		{
			name: "Set UI prompt max rows",
			args: []string{"--ui-prompt-max-rows", "5"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.MaxRows = 5
			}),
		},
		{
			name: "Set UI prompt submit key",
			args: []string{"--ui-prompt-submit-alt=false", "--ui-prompt-submit-key", "space"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.SubmitShortcut = Shortcut{Code: "space"}
			}),
		},
		{
			name: "Set UI prompt pin top",
			args: []string{"--ui-prompt-pin-top=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.PinTop = false
			}),
		},
		{
			name: "Set UI file dialog default directory",
			args: []string{"--ui-file-dialog-default-dir", "/home/user"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.DefaultDirectory = "/home/user"
			}),
		},
		{
			name: "Set UI file dialog show hidden files",
			args: []string{"--ui-file-dialog-show-hidden=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.ShowHiddenFiles = false
			}),
		},
		{
			name: "Set UI file dialog can create directories",
			args: []string{"--ui-file-dialog-can-create-dirs=true"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.CanCreateDirectories = true
			}),
		},
		{
			name: "Set UI file dialog resolves aliases",
			args: []string{"--ui-file-dialog-resolves-aliases=true"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.ResolveAliases = true
			}),
		},
		{
			name: "Set UI file dialog treat packages as directories",
			args: []string{"--ui-file-dialog-treat-packages-as-dirs=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.TreatPackagesAsDirectories = false
			}),
		},
		{
			name: "Set UI file dialog filter display",
			args: []string{"--ui-file-dialog-filter-display=\"Images (*.jpg, *.png)\""},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)"}
			}),
		},
		{
			name: "Set UI file dialog filter pattern",
			args: []string{"--ui-file-dialog-filter-pattern", "*.jpg;*.png"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png"}
			}),
		},
		{
			name: "Set UI file dialog filter display",
			args: []string{"--ui-file-dialog-filter-display=\"Images (*.jpg, *.png)\"", "--ui-file-dialog-filter-display=\"Documents (*.doc, *.docx)\""},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)", "Documents (*.doc, *.docx)"}
			}),
		},
		{
			name: "Set UI file dialog filter pattern",
			args: []string{"--ui-file-dialog-filter-pattern", "*.jpg;*.png", "--ui-file-dialog-filter-pattern", "*.doc;*.docx"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png", "*.doc;*.docx"}
			}),
		},
		{
			name: "Enable UI stream",
			args: []string{"--ui-stream"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Stream = true
			}),
		},
		{
			name: "Enable UI stream - shorthand",
			args: []string{"-s"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Stream = true
			}),
		},
		{
			name: "Set UI title",
			args: []string{"--ui-window-title", "Test Title"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Title = "Test Title"
			}),
		},
		{
			name: "Set UI initial width",
			args: []string{"--ui-window-init-width", "100"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialWidth = expression.NumberContainer{Expression: "100"}
			}),
		},
		{
			name: "Set UI max height",
			args: []string{"--ui-window-max-height", "200"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.MaxHeight = expression.NumberContainer{Expression: "200"}
			}),
		},
		{
			name: "Set UI initial position X",
			args: []string{"--ui-window-init-pos-x", "50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionX = expression.NumberContainer{Expression: "50"}
			}),
		},
		{
			name: "Set UI initial position Y",
			args: []string{"--ui-window-init-pos-y", "50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionY = expression.NumberContainer{Expression: "50"}
			}),
		},
		{
			name: "Set UI initial zoom",
			args: []string{"--ui-window-init-zoom", "1.5"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialZoom = expression.NumberContainer{Expression: "1.5"}
			}),
		},
		{
			name: "Set UI background color",
			args: []string{
				"--ui-window-bg-color-r", "100",
				"--ui-window-bg-color-g", "100",
				"--ui-window-bg-color-b", "100",
				"--ui-window-bg-color-a", "100",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.BackgroundColor = WindowBackgroundColor{R: 100, G: 100, B: 100, A: 100}
			}),
		},
		{
			name: "Set UI start state",
			args: []string{"--ui-window-start-state", "1"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.StartState = 1
			}),
		},
		{
			name: "Disable UI frameless",
			args: []string{"--ui-window-frameless=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Frameless = false
			}),
		},
		{
			name: "Disable UI always on top",
			args: []string{"--ui-window-always-on-top=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.AlwaysOnTop = false
			}),
		},
		{
			name: "Disable UI resizable",
			args: []string{"--ui-window-resizeable=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Resizeable = false
			}),
		},
		{
			name: "Set UI translucent",
			args: []string{"--ui-window-translucent", "never"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Translucent = TranslucentNever
			}),
		},
		{
			name: "Set UI quit shortcut",
			args: []string{
				"--ui-quit-key", "q",
				"--ui-quit-ctrl",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.QuitShortcut = Shortcut{Code: "q", Ctrl: true}
			}),
		},
		{
			name: "Set UI theme",
			args: []string{"--ui-theme", "dark"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Theme = ThemeDark
			}),
		},
		{
			name: "Set UI code style",
			args: []string{"--ui-code-style", "monokai"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.CodeStyle = "monokai"
			}),
		},
		{
			name: "Set UI language",
			args: []string{"--ui-lang", "en_US"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Language = "en_US"
			}),
		},
		{
			name: "Set backend",
			args: []string{"--backend", "openai"},
			expected: modifiedConfig(func(c *Config) {
				c.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set backend - shorthand",
			args: []string{"-b", "openai"},
			expected: modifiedConfig(func(c *Config) {
				c.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set print format",
			args: []string{"--print-format", "plain"},
			expected: modifiedConfig(func(c *Config) {
				c.Printer.Format = PrinterFormatPlain
			}),
		},
		{
			name: "Set print format - shorthand",
			args: []string{"-f", "plain"},
			expected: modifiedConfig(func(c *Config) {
				c.Printer.Format = PrinterFormatPlain
			}),
		},
		{
			name: "Enable print version",
			args: []string{"--version"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.PrintVersion = true
			}),
		},
		{
			name: "Enable print version - shorthand",
			args: []string{"-v"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.PrintVersion = true
			}),
		},
		{
			name: "Set environment variable for init prompt",
			env:  []string{EnvironmentPrefix + "UI_PROMPT_VALUE=test"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set environment variable for log level",
			env:  []string{EnvironmentPrefix + "LOG_LEVEL=-4"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.LogLevel = int(slog.LevelDebug)
			}),
		},
		{
			name: "Set environment variable for pprof address",
			env:  []string{EnvironmentPrefix + "PPROF_ADDRESS=:1312"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.PprofAddress = ":1312"
			}),
		},
		{
			name: "Set environment variable for vue dev tools host",
			env:  []string{EnvironmentPrefix + "VUE_DEV_TOOLS_HOST=localhost"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.VueDevTools.Host = "localhost"
			}),
		},
		{
			name: "Set environment variable for vue dev tools port",
			env:  []string{EnvironmentPrefix + "VUE_DEV_TOOLS_PORT=1312"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.VueDevTools.Port = 1312
			}),
		},
		{
			name: "Set environment variable for open inspector on startup",
			env:  []string{EnvironmentPrefix + "WEBKIT_OPEN_INSPECTOR=1"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.WebKit.OpenInspectorOnStartup = true
			}),
		},
		{
			name: "Set environment variable for webkit http server address",
			env:  []string{EnvironmentPrefix + "WEBKIT_HTTP_SERVER=127.0.0.1:5000"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.WebKit.HttpServerAddress = "127.0.0.1:5000"
			}),
		},
		{
			name: "Set environment variable for UI stream",
			env:  []string{EnvironmentPrefix + "UI_STREAM=true"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Stream = true
			}),
		},
		{
			name: "Set environment variable for UI window title",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_TITLE=Test Title"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Title = "Test Title"
			}),
		},
		{
			name: "Set environment variable for backend",
			env:  []string{EnvironmentPrefix + "BACKEND=openai"},
			expected: modifiedConfig(func(c *Config) {
				c.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set environment variable for print format",
			env:  []string{EnvironmentPrefix + "PRINT_FORMAT=plain"},
			expected: modifiedConfig(func(c *Config) {
				c.Printer.Format = PrinterFormatPlain
			}),
		},
		{
			name: "Set environment variable for UI language",
			env:  []string{EnvironmentPrefix + "UI_LANG=en_US"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Language = "en_US"
			}),
		},
		{
			name: "Set environment variable for UI prompt initial attachments",
			env: []string{
				EnvironmentPrefix + "UI_PROMPT_ATTACHMENTS_1=file2.txt",
				EnvironmentPrefix + "UI_PROMPT_ATTACHMENTS_0=file1.txt",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set environment variable for UI prompt min rows",
			env:  []string{EnvironmentPrefix + "UI_PROMPT_MIN_ROWS=2"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.MinRows = 2
			}),
		},
		{
			name: "Set environment variable for UI prompt max rows",
			env:  []string{EnvironmentPrefix + "UI_PROMPT_MAX_ROWS=5"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.MaxRows = 5
			}),
		},
		{
			name: "Set environment variable for UI prompt submit key",
			env: []string{
				EnvironmentPrefix + "UI_PROMPT_SUBMIT_ALT=false",
				EnvironmentPrefix + "UI_PROMPT_SUBMIT_KEY=space",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.SubmitShortcut = Shortcut{Code: "space"}
			}),
		},
		{
			name: "Set environment variable for UI prompt pin top",
			env: []string{
				EnvironmentPrefix + "UI_PROMPT_PIN_TOP=false",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.PinTop = false
			}),
		},
		{
			name: "Set environment variable for UI file dialog default directory",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_DEFAULT_DIR=/home/user"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.DefaultDirectory = "/home/user"
			}),
		},
		{
			name: "Set environment variable for UI file dialog show hidden files",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_SHOW_HIDDEN=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.ShowHiddenFiles = false
			}),
		},
		{
			name: "Set environment variable for UI file dialog can create directories",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_CAN_CREATE_DIRS=true"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.CanCreateDirectories = true
			}),
		},
		{
			name: "Set environment variable for UI file dialog resolves aliases",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_RESOLVES_ALIASES=true"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.ResolveAliases = true
			}),
		},
		{
			name: "Set environment variable for UI file dialog treat packages as directories",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_TREAT_PACKAGES_AS_DIRS=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.TreatPackagesAsDirectories = false
			}),
		},
		{
			name: "Set environment variable for UI file dialog filter display",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_FILTER_DISPLAY_0=Images (*.jpg, *.png)"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)"}
			}),
		},
		{
			name: "Set environment variable for UI file dialog filter pattern",
			env:  []string{EnvironmentPrefix + "UI_FILE_DIALOG_FILTER_PATTERN_0=*.jpg;*.png"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png"}
			}),
		},
		{
			name: "Set environment variable for UI initial width",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_INIT_WIDTH=100"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialWidth = expression.NumberContainer{Expression: "100"}
			}),
		},
		{
			name: "Set environment variable for UI max height",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_MAX_HEIGHT=200"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.MaxHeight = expression.NumberContainer{Expression: "200"}
			}),
		},
		{
			name: "Set environment variable for UI initial position X",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_INIT_POS_X=50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionX = expression.NumberContainer{Expression: "50"}
			}),
		},
		{
			name: "Set environment variable for UI initial position Y",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_INIT_POS_Y=50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionY = expression.NumberContainer{Expression: "50"}
			}),
		},
		{
			name: "Set environment variable for UI initial zoom",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_INIT_ZOOM=1.5"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialZoom = expression.NumberContainer{Expression: "1.5"}
			}),
		},
		{
			name: "Set environment variable for UI background color",
			env: []string{
				EnvironmentPrefix + "UI_WINDOW_BG_COLOR_R=100",
				EnvironmentPrefix + "UI_WINDOW_BG_COLOR_G=100",
				EnvironmentPrefix + "UI_WINDOW_BG_COLOR_B=100",
				EnvironmentPrefix + "UI_WINDOW_BG_COLOR_A=100",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.BackgroundColor = WindowBackgroundColor{R: 100, G: 100, B: 100, A: 100}
			}),
		},
		{
			name: "Set environment variable for UI start state",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_START_STATE=1"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.StartState = 1
			}),
		},
		{
			name: "Set environment variable for UI frameless",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_FRAMELESS=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Frameless = false
			}),
		},
		{
			name: "Set environment variable for UI always on top",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_ALWAYS_ON_TOP=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.AlwaysOnTop = false
			}),
		},
		{
			name: "Set environment variable for UI resizable",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_RESIZEABLE=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Resizeable = false
			}),
		},
		{
			name: "Set environment variable for UI translucent",
			env:  []string{EnvironmentPrefix + "UI_WINDOW_TRANSLUCENT=never"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Translucent = TranslucentNever
			}),
		},
		{
			name: "Set environment variable for UI quit shortcut",
			env: []string{
				EnvironmentPrefix + "UI_QUIT_KEY=q",
				EnvironmentPrefix + "UI_QUIT_CTRL=true",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.QuitShortcut = Shortcut{Code: "q", Ctrl: true}
			}),
		},
		{
			name: "Set environment variable for UI theme",
			env:  []string{EnvironmentPrefix + "UI_THEME=dark"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Theme = ThemeDark
			}),
		},
		{
			name: "Set environment variable for UI code style",
			env:  []string{EnvironmentPrefix + "UI_CODE_STYLE=monokai"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.CodeStyle = "monokai"
			}),
		},
		{
			name: "Set environment variable for print version",
			env:  []string{EnvironmentPrefix + "VERSION=true"},
			expected: modifiedConfig(func(c *Config) {
				c.Debug.PrintVersion = true
			}),
		},
		{
			name: "Argument will override environment",
			args: []string{"--ui-prompt-value=arg-prompt"},
			env:  []string{EnvironmentPrefix + "UI_PROMPT_VALUE=env-prompt"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Prompt.InitValue = "arg-prompt"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Parse(tt.args, tt.env)
			assert.Equal(t, tt.expected, *c)
		})
	}
}
