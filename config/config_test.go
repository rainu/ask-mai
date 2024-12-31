package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v2/pkg/options"
	"io"
	"log/slog"
	"os"
	"testing"
)

func defaultConfig() Config {
	return Config{
		LogLevel: int(slog.LevelError),
		UI: UIConfig{
			Prompt: PromptConfig{
				InitValue:       "",
				InitAttachments: []string{},
				MinRows:         1,
				MaxRows:         4,
				SubmitShortcut:  Shortcut{Alt: true, Code: "enter"},
			},
			FileDialog: FileDialogConfig{
				DefaultDirectory:           "",
				ShowHiddenFiles:            true,
				CanCreateDirectories:       false,
				ResolvesAliases:            true,
				TreatPackagesAsDirectories: true,
				FilterDisplay:              []string{},
				FilterPattern:              []string{},
			},
			Stream: false,
			Window: WindowConfig{
				Title:            "Prompt - Ask mAI",
				InitialWidth:     ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Width/2"},
				MaxHeight:        ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Height/3"},
				InitialPositionX: ExpressionContainer{Expression: "v.CurrentScreen.Dimension.Width/4"},
				InitialPositionY: ExpressionContainer{Expression: "0"},
				InitialZoom:      ExpressionContainer{Expression: "1.0"},
				BackgroundColor:  WindowBackgroundColor{R: 255, G: 255, B: 255, A: 192},
				StartState:       int(options.Normal),
				Frameless:        true,
				AlwaysOnTop:      true,
				Resizeable:       true,
				Translucent:      TranslucentHover,
			},
			QuitShortcut: Shortcut{Code: "escape"},
			Theme:        ThemeSystem,
			CodeStyle:    "github",
			Language:     os.Getenv("LANG"),
		},
		Backend: BackendCopilot,
		OpenAI: OpenAIConfig{
			APIType: "OPEN_AI",
			Model:   "gpt-4o-mini",
		},
		CallOptions: CallOptionsConfig{
			Temperature: -1,
			TopK:        -1,
			TopP:        -1,
		},
		PrintVersion: false,
		Printer: PrinterConfig{
			Format:  PrinterFormatJSON,
			Targets: []io.WriteCloser{os.Stdout},
			targets: PrinterTargetOut,
		},
	}
}

func modifiedConfig(mod func(*Config)) Config {
	c := defaultConfig()
	mod(&c)
	return c
}

func TestConfig_Parse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name:     "Default values",
			args:     []string{},
			expected: defaultConfig(),
		},
		{
			name: "Set log level",
			args: []string{"--log-level", "-4"},
			expected: modifiedConfig(func(c *Config) {
				c.LogLevel = int(slog.LevelDebug)
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
				c.CallOptions.SystemPrompt = "test"
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
				c.UI.FileDialog.ResolvesAliases = true
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
			args: []string{"--ui-title", "Test Title"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Title = "Test Title"
			}),
		},
		{
			name: "Set UI initial width",
			args: []string{"--ui-init-width", "100"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialWidth = ExpressionContainer{Expression: "100"}
			}),
		},
		{
			name: "Set UI max height",
			args: []string{"--ui-max-height", "200"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.MaxHeight = ExpressionContainer{Expression: "200"}
			}),
		},
		{
			name: "Set UI initial position X",
			args: []string{"--ui-init-pos-x", "50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionX = ExpressionContainer{Expression: "50"}
			}),
		},
		{
			name: "Set UI initial position Y",
			args: []string{"--ui-init-pos-y", "50"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialPositionY = ExpressionContainer{Expression: "50"}
			}),
		},
		{
			name: "Set UI initial zoom",
			args: []string{"--ui-init-zoom", "1.5"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.InitialZoom = ExpressionContainer{Expression: "1.5"}
			}),
		},
		{
			name: "Set UI background color",
			args: []string{
				"--ui-bg-color-r", "100",
				"--ui-bg-color-g", "100",
				"--ui-bg-color-b", "100",
				"--ui-bg-color-a", "100",
			},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.BackgroundColor = WindowBackgroundColor{R: 100, G: 100, B: 100, A: 100}
			}),
		},
		{
			name: "Set UI start state",
			args: []string{"--ui-start-state", "1"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.StartState = 1
			}),
		},
		{
			name: "Disable UI frameless",
			args: []string{"--ui-frameless=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Frameless = false
			}),
		},
		{
			name: "Disable UI always on top",
			args: []string{"--ui-always-on-top=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.AlwaysOnTop = false
			}),
		},
		{
			name: "Disable UI resizable",
			args: []string{"--ui-resizeable=false"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Resizeable = false
			}),
		},
		{
			name: "Set UI translucent",
			args: []string{"--ui-translucent", "never"},
			expected: modifiedConfig(func(c *Config) {
				c.UI.Window.Translucent = TranslucentNever
			}),
		},
		{
			name: "Set UI quit shortcut",
			args: []string{
				"--ui-quit-shortcut-keycode", "q",
				"--ui-quit-shortcut-ctrl",
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
				c.Backend = BackendOpenAI
			}),
		},
		{
			name: "Set backend - shorthand",
			args: []string{"-b", "openai"},
			expected: modifiedConfig(func(c *Config) {
				c.Backend = BackendOpenAI
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
			name: "Set print targets",
			args: []string{"--print-targets", "stderr"},
			expected: modifiedConfig(func(c *Config) {
				c.Printer.targets = PrinterTargetErr
			}),
		},
		{
			name: "Set print targets - shorthand",
			args: []string{"-o", "stderr"},
			expected: modifiedConfig(func(c *Config) {
				c.Printer.targets = PrinterTargetErr
			}),
		},
		{
			name: "Enable print version",
			args: []string{"--version"},
			expected: modifiedConfig(func(c *Config) {
				c.PrintVersion = true
			}),
		},
		{
			name: "Enable print version - shorthand",
			args: []string{"-v"},
			expected: modifiedConfig(func(c *Config) {
				c.PrintVersion = true
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Parse(tt.args)
			assert.Equal(t, tt.expected, *c)
		})
	}
}
