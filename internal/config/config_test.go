package config

import (
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/config/model/common"
	"github.com/rainu/ask-mai/internal/config/model/llm"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/mcp"
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func modifiedConfig(mod func(*model.Config)) model.Config {
	c := &model.Config{}
	yacl.NewConfig(c).ApplyDefaults()
	mod(c)
	return *c
}

func TestConfig_Parse(t *testing.T) {
	origFn := yamlLookupLocations
	defer func() {
		yamlLookupLocations = origFn
	}()
	yamlLookupLocations = func() (result []string) {
		return []string{}
	}

	tests := []struct {
		name     string
		args     []string
		env      []string
		expected model.Config
	}{
		{
			name:     "Default values",
			args:     []string{},
			expected: modifiedConfig(func(c *model.Config) {}),
		},
		{
			name: "Set log level",
			args: []string{"--log-level=debug"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.LogLevel = model.LogLevelDebug
			}),
		},
		{
			name: "Set pprof address",
			args: []string{"--pprof-address=localhost:6060"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.PprofAddress = "localhost:6060"
			}),
		},
		{
			name: "Set vue dev tools host",
			args: []string{"--vue-dev-tools.host=localhost"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.VueDevTools.Host = "localhost"
			}),
		},
		{
			name: "Set vue dev tools port",
			args: []string{"--vue-dev-tools.port=1312"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.VueDevTools.Port = 1312
			}),
		},
		{
			name: "Set open inspector on startup",
			args: []string{"--webkit.open-inspector"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.WebKit.OpenInspectorOnStartup = true
			}),
		},
		{
			name: "Set webkit inspector http server address",
			args: []string{"--webkit.http-server=127.0.0.1:5000"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.WebKit.HttpServerAddress = "127.0.0.1:5000"
			}),
		},
		{
			name: "Set UI prompt value",
			args: []string{"--llm.call.prompt.init-value=test"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set UI prompt value - shorthand",
			args: []string{"-p=test"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set system prompt value",
			args: []string{"--llm.call.prompt.system=test"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.System = "test"
			}),
		},
		{
			name: "Set initial tool call arguments",
			args: []string{
				"--llm.call.prompt.init-tool-call[0].server=_builtin",
				"--llm.call.prompt.init-tool-call[0].name=test",
				"--llm.call.prompt.init-tool-call[0].args.key=value",
				"--llm.call.prompt.init-tool-call[1].server=_builtin",
				"--llm.call.prompt.init-tool-call[1].name=test",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitToolCalls = []llm.ToolCall{
					{
						Server: "_builtin",
						Name:   "test",
						Arguments: map[string]any{
							"key": "value",
						},
					},
					{
						Server: "_builtin",
						Name:   "test",
					},
				}
			}),
		},
		{
			name: "Set UI prompt initial attachments",
			args: []string{"--llm.call.prompt.init-attachment.[0]=file1.txt", "--llm.call.prompt.init-attachment.[1]=file2.txt"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set UI prompt initial attachments - shorthand",
			args: []string{"-a=file1.txt", "-a=file2.txt"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set UI prompt min rows",
			args: []string{"--ui.prompt.min-rows=2"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.MinRows = yacl.P(uint(2))
			}),
		},
		{
			name: "Set UI prompt max rows",
			args: []string{"--ui.prompt.max-rows=5"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.MaxRows = yacl.P(uint(5))
			}),
		},
		{
			name: "Set UI prompt submit key",
			args: []string{
				"--ui.prompt.submit.binding.[0]=space",
				"--ui.prompt.submit.binding.[1]=escape",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.SubmitShortcut = model.Shortcut{Binding: []string{"space", "escape"}}
			}),
		},
		{
			name: "Set UI prompt pin top",
			args: []string{"--ui.prompt.pin-top=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.PinTop = yacl.P(false)
			}),
		},
		{
			name: "Set UI file dialog default directory",
			args: []string{"--ui.file-dialog.default-dir=/home/user"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.DefaultDirectory = "/home/user"
			}),
		},
		{
			name: "Set UI file dialog show hidden files",
			args: []string{"--ui.file-dialog.show-hidden=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.ShowHiddenFiles = yacl.P(false)
			}),
		},
		{
			name: "Set UI file dialog can create directories",
			args: []string{"--ui.file-dialog.can-create-dirs=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.CanCreateDirectories = yacl.P(true)
			}),
		},
		{
			name: "Set UI file dialog resolves aliases",
			args: []string{"--ui.file-dialog.resolve-aliases=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.ResolveAliases = yacl.P(true)
			}),
		},
		{
			name: "Set UI file dialog treat packages as directories",
			args: []string{"--ui.file-dialog.treat-packages-as-dirs=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.TreatPackagesAsDirectories = yacl.P(false)
			}),
		},
		{
			name: "Set UI file dialog filter display",
			args: []string{"--ui.file-dialog.filter-display.[0]=Images (*.jpg, *.png)"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)"}
			}),
		},
		{
			name: "Set UI file dialog filter pattern",
			args: []string{"--ui.file-dialog.filter-pattern=*.jpg;*.png"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png"}
			}),
		},
		{
			name: "Set UI file dialog filter display",
			args: []string{"--ui.file-dialog.filter-display=Images (*.jpg, *.png)", "--ui.file-dialog.filter-display=Documents (*.doc, *.docx)"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)", "Documents (*.doc, *.docx)"}
			}),
		},
		{
			name: "Set UI file dialog filter pattern",
			args: []string{"--ui.file-dialog.filter-pattern=*.jpg;*.png", "--ui.file-dialog.filter-pattern=*.doc;*.docx"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png", "*.doc;*.docx"}
			}),
		},
		{
			name: "Enable UI stream",
			args: []string{"--ui.stream"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Stream = yacl.P(true)
			}),
		},
		{
			name: "Enable UI stream - shorthand",
			args: []string{"-s"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Stream = yacl.P(true)
			}),
		},
		{
			name: "Set UI title",
			args: []string{"--ui.window.title=Test Title"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Title = "Test Title"
			}),
		},
		{
			name: "Set UI initial width",
			args: []string{"--ui.window.init-width.expression=100"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialWidth.Expression = yacl.P("100")
			}),
		},
		{
			name: "Set UI max height",
			args: []string{"--ui.window.max-height.expression=200"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.MaxHeight = common.NumberContainer{Expression: yacl.P("200")}
			}),
		},
		{
			name: "Set UI initial position X",
			args: []string{"--ui.window.init-pos-x.expression=50"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialPositionX = common.NumberContainer{Expression: yacl.P("50")}
			}),
		},
		{
			name: "Set UI initial position Y",
			args: []string{"--ui.window.init-pos-y.expression=50"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialPositionY = common.NumberContainer{Expression: yacl.P("50")}
			}),
		},
		{
			name: "Set UI initial zoom",
			args: []string{"--ui.window.init-zoom.expression=1.5"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialZoom = common.NumberContainer{Expression: yacl.P("1.5")}
			}),
		},
		{
			name: "Set UI background color",
			args: []string{
				"--ui.window.bg-color.r=100",
				"--ui.window.bg-color.g=100",
				"--ui.window.bg-color.b=100",
				"--ui.window.bg-color.a=100",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.BackgroundColor = model.WindowBackgroundColor{R: yacl.P(uint(100)), G: yacl.P(uint(100)), B: yacl.P(uint(100)), A: yacl.P(uint(100))}
			}),
		},
		{
			name: "Set UI start state",
			args: []string{"--ui.window.start-state=1"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.StartState = yacl.P(1)
			}),
		},
		{
			name: "Disable UI frameless",
			args: []string{"--ui.window.frameless=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Frameless = yacl.P(false)
			}),
		},
		{
			name: "Disable UI always on top",
			args: []string{"--ui.window.always-on-top=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.AlwaysOnTop = yacl.P(false)
			}),
		},
		{
			name: "Disable UI resizable",
			args: []string{"--ui.window.resizeable=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Resizeable = yacl.P(false)
			}),
		},
		{
			name: "Set UI translucent",
			args: []string{"--ui.window.translucent=never"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Translucent = model.TranslucentNever
			}),
		},
		{
			name: "Set UI quit shortcut",
			args: []string{
				"--ui.quit.binding=ctrl+q",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.QuitShortcut = model.Shortcut{Binding: []string{"ctrl+q"}}
			}),
		},
		{
			name: "Set UI theme",
			args: []string{"--ui.theme=dark"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Theme = model.ThemeDark
			}),
		},
		{
			name: "Set UI code style",
			args: []string{"--ui.code-style=monokai"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.CodeStyle = "monokai"
			}),
		},
		{
			name: "Set UI language",
			args: []string{"--ui.lang=en_US"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Language = "en_US"
			}),
		},
		{
			name: "Set backend",
			args: []string{"--llm.backend=openai"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set backend - shorthand",
			args: []string{"-b=openai"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set print format",
			args: []string{"--print.format=plain"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.Printer.Format = model.PrinterFormatPlain
			}),
		},
		{
			name: "Set print format - shorthand",
			args: []string{"-f=plain"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.Printer.Format = model.PrinterFormatPlain
			}),
		},
		{
			name: "Enable print version",
			args: []string{"--version"},
			expected: modifiedConfig(func(c *model.Config) {
				c.Version = true
			}),
		},
		{
			name: "Enable print version - shorthand",
			args: []string{"-v"},
			expected: modifiedConfig(func(c *model.Config) {
				c.Version = true
			}),
		},
		{
			name: "Set environment variable for init prompt",
			env:  []string{EnvironmentPrefix + "=--llm.call.prompt.init-value=test"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitValue = "test"
			}),
		},
		{
			name: "Set environment variable for log level",
			env:  []string{EnvironmentPrefix + "=--log-level=debug"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.LogLevel = model.LogLevelDebug
			}),
		},
		{
			name: "Set environment variable for pprof address",
			env:  []string{EnvironmentPrefix + "=--pprof-address=:1312"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.PprofAddress = ":1312"
			}),
		},
		{
			name: "Set environment variable for vue dev tools host",
			env:  []string{EnvironmentPrefix + "=--vue-dev-tools.host=localhost"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.VueDevTools.Host = "localhost"
			}),
		},
		{
			name: "Set environment variable for vue dev tools port",
			env:  []string{EnvironmentPrefix + "=--vue-dev-tools.port=1312"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.VueDevTools.Port = 1312
			}),
		},
		{
			name: "Set environment variable for open inspector on startup",
			env:  []string{EnvironmentPrefix + "=--webkit.open-inspector"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.WebKit.OpenInspectorOnStartup = true
			}),
		},
		{
			name: "Set environment variable for webkit http server address",
			env:  []string{EnvironmentPrefix + "=--webkit.http-server=127.0.0.1:5000"},
			expected: modifiedConfig(func(c *model.Config) {
				c.DebugConfig.WebKit.HttpServerAddress = "127.0.0.1:5000"
			}),
		},
		{
			name: "Set environment variable for UI stream",
			env:  []string{EnvironmentPrefix + "=--ui.stream=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Stream = yacl.P(true)
			}),
		},
		{
			name: "Set environment variable for UI window title",
			env:  []string{EnvironmentPrefix + "=--ui.window.title=Test Title"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Title = "Test Title"
			}),
		},
		{
			name: "Set environment variable for backend",
			env:  []string{EnvironmentPrefix + "=--llm.backend=openai"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.Backend = "openai"
			}),
		},
		{
			name: "Set environment variable for print format",
			env:  []string{EnvironmentPrefix + "=--print.format=plain"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.Printer.Format = model.PrinterFormatPlain
			}),
		},
		{
			name: "Set environment variable for UI language",
			env:  []string{EnvironmentPrefix + "=--ui.lang=en_US"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Language = "en_US"
			}),
		},
		{
			name: "Set environment variable for UI prompt initial attachments",
			env: []string{
				EnvironmentPrefix + "=--llm.call.prompt.init-attachment.[1]=file2.txt",
				EnvironmentPrefix + "=--llm.call.prompt.init-attachment.[0]=file1.txt",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitAttachments = []string{"file1.txt", "file2.txt"}
			}),
		},
		{
			name: "Set environment variable for UI prompt min rows",
			env:  []string{EnvironmentPrefix + "=--ui.prompt.min-rows=2"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.MinRows = yacl.P(uint(2))
			}),
		},
		{
			name: "Set environment variable for UI prompt max rows",
			env:  []string{EnvironmentPrefix + "=--ui.prompt.max-rows=5"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.MaxRows = yacl.P(uint(5))
			}),
		},
		{
			name: "Set environment variable for UI prompt submit key",
			env: []string{
				EnvironmentPrefix + "=--ui.prompt.submit.binding=space",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.SubmitShortcut = model.Shortcut{Binding: []string{"space"}}
			}),
		},
		{
			name: "Set environment variable for UI prompt pin top",
			env: []string{
				EnvironmentPrefix + "=--ui.prompt.pin-top=false",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Prompt.PinTop = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI file dialog default directory",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.default-dir=/home/user"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.DefaultDirectory = "/home/user"
			}),
		},
		{
			name: "Set environment variable for UI file dialog show hidden files",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.show-hidden=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.ShowHiddenFiles = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI file dialog can create directories",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.can-create-dirs=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.CanCreateDirectories = yacl.P(true)
			}),
		},
		{
			name: "Set environment variable for UI file dialog resolves aliases",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.resolve-aliases=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.ResolveAliases = yacl.P(true)
			}),
		},
		{
			name: "Set environment variable for UI file dialog treat packages as directories",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.treat-packages-as-dirs=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.TreatPackagesAsDirectories = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI file dialog filter display",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.filter-display=Images (*.jpg, *.png)"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterDisplay = []string{"Images (*.jpg, *.png)"}
			}),
		},
		{
			name: "Set environment variable for UI file dialog filter pattern",
			env:  []string{EnvironmentPrefix + "=--ui.file-dialog.filter-pattern=*.jpg;*.png"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.FileDialog.FilterPattern = []string{"*.jpg;*.png"}
			}),
		},
		{
			name: "Set environment variable for UI initial width",
			env:  []string{EnvironmentPrefix + "=--ui.window.init-width.expression=100"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialWidth = common.NumberContainer{Expression: yacl.P("100")}
			}),
		},
		{
			name: "Set environment variable for UI max height",
			env:  []string{EnvironmentPrefix + "=--ui.window.max-height.expression=200"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.MaxHeight = common.NumberContainer{Expression: yacl.P("200")}
			}),
		},
		{
			name: "Set environment variable for UI initial position X",
			env:  []string{EnvironmentPrefix + "=--ui.window.init-pos-x.expression=50"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialPositionX = common.NumberContainer{Expression: yacl.P("50")}
			}),
		},
		{
			name: "Set environment variable for UI initial position Y",
			env:  []string{EnvironmentPrefix + "=--ui.window.init-pos-y.expression=50"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialPositionY = common.NumberContainer{Expression: yacl.P("50")}
			}),
		},
		{
			name: "Set environment variable for UI initial zoom",
			env:  []string{EnvironmentPrefix + "=--ui.window.init-zoom.expression=1.5"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.InitialZoom = common.NumberContainer{Expression: yacl.P("1.5")}
			}),
		},
		{
			name: "Set environment variable for UI background color",
			env: []string{
				EnvironmentPrefix + "=--ui.window.bg-color.r=100",
				EnvironmentPrefix + "=--ui.window.bg-color.g=100",
				EnvironmentPrefix + "=--ui.window.bg-color.b=100",
				EnvironmentPrefix + "=--ui.window.bg-color.a=100",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.BackgroundColor = model.WindowBackgroundColor{R: yacl.P(uint(100)), G: yacl.P(uint(100)), B: yacl.P(uint(100)), A: yacl.P(uint(100))}
			}),
		},
		{
			name: "Set environment variable for UI start state",
			env:  []string{EnvironmentPrefix + "=--ui.window.start-state=1"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.StartState = yacl.P(1)
			}),
		},
		{
			name: "Set environment variable for UI frameless",
			env:  []string{EnvironmentPrefix + "=--ui.window.frameless=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Frameless = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI always on top",
			env:  []string{EnvironmentPrefix + "=--ui.window.always-on-top=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.AlwaysOnTop = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI resizable",
			env:  []string{EnvironmentPrefix + "=--ui.window.resizeable=false"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Resizeable = yacl.P(false)
			}),
		},
		{
			name: "Set environment variable for UI translucent",
			env:  []string{EnvironmentPrefix + "=--ui.window.translucent=never"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Window.Translucent = model.TranslucentNever
			}),
		},
		{
			name: "Set environment variable for UI quit shortcut",
			env: []string{
				EnvironmentPrefix + "=--ui.quit.binding=ctrl+q",
			},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.QuitShortcut = model.Shortcut{Binding: []string{"ctrl+q"}}
			}),
		},
		{
			name: "Set environment variable for UI theme",
			env:  []string{EnvironmentPrefix + "=--ui.theme=dark"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.Theme = model.ThemeDark
			}),
		},
		{
			name: "Set environment variable for UI code style",
			env:  []string{EnvironmentPrefix + "=--ui.code-style=monokai"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.UI.CodeStyle = "monokai"
			}),
		},
		{
			name: "Set environment variable for print version",
			env:  []string{EnvironmentPrefix + "=--version=true"},
			expected: modifiedConfig(func(c *model.Config) {
				c.Version = true
			}),
		},
		{
			name: "Argument will override environment",
			args: []string{"--llm.call.prompt.init-value=arg-prompt"},
			env:  []string{EnvironmentPrefix + "=--llm.call.prompt.init-value=env-prompt"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.CallOptions.Prompt.InitValue = "arg-prompt"
			}),
		},
		{
			name: "Set Secret - Command",
			args: []string{"--llm.openai.api-key.command.name=echo", "--llm.openai.api-key.command.args=secret", "--llm.openai.api-key.command.no-trim"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.OpenAI.APIKey.Command.Name = "echo"
				c.MainProfile.LLM.OpenAI.APIKey.Command.Args = []string{"secret"}
				c.MainProfile.LLM.OpenAI.APIKey.Command.NoTrim = true
			}),
		},
		{
			name: "Set Secret - Plain",
			args: []string{"--llm.openai.api-key.plain=secret"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.OpenAI.APIKey.Plain = "secret"
			}),
		},
		{
			name: "Set MCP Timeout",
			args: []string{"--llm.tool.mcpServers[test].timeout.execution=1m"},
			expected: modifiedConfig(func(c *model.Config) {
				c.MainProfile.LLM.Tool.McpServer = map[string]mcp.Server{
					"test": {
						Timeout: mcp.Timeout{
							Init:      yacl.P(5 * time.Second),
							List:      yacl.P(5 * time.Second),
							Execution: yacl.P(1 * time.Minute),
						},
					},
				}
			}),
		},
		{
			name: "Prompt with another active-profile",
			args: []string{"--profiles.test.description=test", "-P", "test", "-p", "What is the answer?"},
			expected: modifiedConfig(func(c *model.Config) {
				c.ActiveProfile = "test"
				c.MainProfile.LLM.CallOptions.Prompt.InitValue = "What is the answer?"

				testProfile := &model.Profile{}
				yacl.NewConfig(testProfile).ApplyDefaults()
				testProfile.Meta.Description = "test"
				testProfile.LLM.CallOptions.Prompt.InitValue = "What is the answer?"

				c.Profiles = map[string]*model.Profile{
					"test": testProfile,
				}
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Parse(tt.args, tt.env)
			c.MainProfile.Printer.Targets = nil
			assert.Equal(t, tt.expected, *c)
		})
	}
}
