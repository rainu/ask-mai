package config

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/config/model/common"
	"github.com/rainu/ask-mai/internal/config/model/llm"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	iMcp "github.com/rainu/ask-mai/internal/config/model/llm/tools/mcp"
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func Test_processYaml(t *testing.T) {
	c := &model.Config{}

	yamlContent := `
ui:
  window:
    title: Test Window
    init-width:
      expression: "800"
      value: 0
    max-height:
      expression: "600"
      value: 0
    init-pos-x:
      expression: "100"
      value: 0
    init-pos-y:
      expression: "100"
      value: 0
    init-zoom:
      expression: "1.0"
      value: 0
    bg-color:
      r: 255
      g: 255
      b: 255
      a: 255
    start-state: 1
    always-on-top: true
    frameless: true
    resizeable: true
    translucent: never
  prompt:
    value: Initial Prompt
    attachments:
      - attachment1
      - attachment2
    min-rows: 1
    max-rows: 10
    submit:
      binding:
        - "alt+ctrl+meta+shift+enter"
  file-dialog:
    default-dir: /root
    show-hidden: true
    can-create-dirs: true
    resolve-aliases: true
    treat-packages-as-dirs: true
    filter-display:
      - Image
    filter-pattern:
      - '*.png'
  stream: true
  quit:
    binding:
      - "alt+ctrl+meta+shift+escape"
  theme: dark
  code-style: default
  lang: en
llm:
  backend: anthropic
  localai:
    api-key: 
      plain: APIKey
    model: model
    base-url: baseurl
  openai:
    api-key:
      command: 
        name: echo
        args:
          - APIKey
    api-type: APIType
    api-version: APIVersion
    model: Model
    base-url: BaseUrl
    organization: Organization
  anythingllm:
    base-url: BaseURL
    token:
      command: 
        name: echo
        args:
          - '-n'
          - Token
        no-trim: true
    workspace: Workspace
  ollama:
    server-url: ServerURL
    model: Model
  mistral:
    api-key:
      command: 
        name: echo
        args:
          - '-n'
          - ApiKey
    endpoint: Endpoint
    model: Model
  anthropic:
    api-key: 
      plain: Token
    base-url: BaseUrl
    model: Model
  call:
    system-prompt: Your system prompt
    max-token: 1000
    temperature: 0.7
    top-k: 50
    top-p: 0.9
    min-length: 10
    max-length: 200
  tool:
    builtin:
      command-execution:
        disable: true
    custom:
      test:
        description: This is a test function.
        parameters:
          type: object
          properties:
            arg1:
              type: string
              description: The first argument.
            arg2:
              type: number
              description: The second argument.
          required:
            - arg1
        command: doTest.sh
        "approval": true
    mcpServers:
      command1:
        command: docker
        args:
          - run
          - --rm
          - -i
          - -e
          - GITHUB_PERSONAL_ACCESS_TOKEN=github_
          - ghcr.io/github/github-mcp-server
        env:
          TEST: test
        timeout:
          init: 500ms
          list: 10s
          execution: 1m30s
      command2:
        command: echo
      http:
        baseUrl: http://localhost:8080/api/v1
        headers:
          Authorization: Bearer TOKEN
print:
  format: json
  targets:
    - stdout
log-level: debug
pprof-address: ":1312"
vue-dev-tools:
  host: "localhost"
  port: 1312
webkit:
  open-inspector: true
  http-server: "127.0.0.1:5000"
themes:
  dark:
    colors:
      background: "#FFFFFF"
      surface: "#FFFFFF"
  light:
    colors:
      chat-tool-call: "#FF0000"
  custom:
    test:
      colors:
         background: "#00FF00"
         chat-tool-call: "#0000FF"
`
	// add profile "test" with the same values as default
	yamlContent += "\nprofiles:\n  test:\n" + strings.ReplaceAll(yamlContent, "\n", "\n    ")

	sr := strings.NewReader(yamlContent)
	config := yacl.NewConfig(c, yacl.WithAutoApplyDefaults(false))

	require.NoError(t, processYaml(config, sr))

	expDefConf := model.Config{
		MainProfile: model.Profile{
			UI: model.UIConfig{
				Window: model.WindowConfig{
					Title:            "Test Window",
					InitialWidth:     common.NumberContainer{Expression: yacl.P("800"), Value: yacl.P(0.0)},
					MaxHeight:        common.NumberContainer{Expression: yacl.P("600"), Value: yacl.P(0.0)},
					InitialPositionX: common.NumberContainer{Expression: yacl.P("100"), Value: yacl.P(0.0)},
					InitialPositionY: common.NumberContainer{Expression: yacl.P("100"), Value: yacl.P(0.0)},
					InitialZoom:      common.NumberContainer{Expression: yacl.P("1.0"), Value: yacl.P(0.0)},
					BackgroundColor:  model.WindowBackgroundColor{R: yacl.P(uint(255)), G: yacl.P(uint(255)), B: yacl.P(uint(255)), A: yacl.P(uint(255))},
					StartState:       yacl.P(1),
					AlwaysOnTop:      yacl.P(true),
					Frameless:        yacl.P(true),
					Resizeable:       yacl.P(true),
					Translucent:      "never",
				},
				Prompt: model.PromptConfig{
					InitValue:       "Initial Prompt",
					InitAttachments: []string{"attachment1", "attachment2"},
					MinRows:         yacl.P(uint(1)),
					MaxRows:         yacl.P(uint(10)),
					SubmitShortcut: model.Shortcut{
						Binding: []string{"alt+ctrl+meta+shift+enter"},
					},
				},
				FileDialog: model.FileDialogConfig{
					DefaultDirectory:           "/root",
					ShowHiddenFiles:            yacl.P(true),
					CanCreateDirectories:       yacl.P(true),
					ResolveAliases:             yacl.P(true),
					TreatPackagesAsDirectories: yacl.P(true),
					FilterDisplay:              []string{"Image"},
					FilterPattern:              []string{"*.png"},
				},
				Stream: yacl.P(true),
				QuitShortcut: model.Shortcut{
					Binding: []string{"alt+ctrl+meta+shift+escape"},
				},
				Theme:     "dark",
				CodeStyle: "default",
				Language:  "en",
			},
			LLM: llm.LLMConfig{
				Backend: "anthropic",
				LocalAI: llm.LocalAIConfig{
					APIKey: common.Secret{
						Plain: "APIKey",
					},
					Model:   "model",
					BaseUrl: "baseurl",
				},
				OpenAI: llm.OpenAIConfig{
					APIKey: common.Secret{
						Command: common.SecretCommand{
							Name: "echo",
							Args: []string{"APIKey"},
						},
					},
					APIType:      "APIType",
					APIVersion:   "APIVersion",
					Model:        "Model",
					BaseUrl:      "BaseUrl",
					Organization: "Organization",
				},
				AnythingLLM: llm.AnythingLLMConfig{
					BaseURL: "BaseURL",
					Token: common.Secret{
						Command: common.SecretCommand{
							Name:   "echo",
							Args:   []string{"-n", "Token"},
							NoTrim: true,
						},
					},
					Workspace: "Workspace",
				},
				Ollama: llm.OllamaConfig{
					ServerURL: "ServerURL",
					Model:     "Model",
				},
				Mistral: llm.MistralConfig{
					ApiKey: common.Secret{
						Command: common.SecretCommand{
							Name: "echo",
							Args: []string{"-n", "ApiKey"},
						},
					},
					Endpoint: "Endpoint",
					Model:    "Model",
				},
				Anthropic: llm.AnthropicConfig{
					Token: common.Secret{
						Plain: "Token",
					},
					BaseUrl: "BaseUrl",
					Model:   "Model",
				},
				CallOptions: llm.CallOptionsConfig{
					SystemPrompt: "Your system prompt",
					MaxToken:     1000,
					Temperature:  0.7,
					TopK:         50,
					TopP:         0.9,
					MinLength:    10,
					MaxLength:    200,
				},
				Tool: tools.Config{
					BuiltIns: builtin.BuiltIns{
						CommandExec: builtin.CommandExecution{
							Disable: true,
						},
					},
					Custom: map[string]command.FunctionDefinition{
						"test": {
							Description: "This is a test function.",
							Parameters: mcp.ToolInputSchema{
								Type: "object",
								Properties: map[string]any{
									"arg1": map[string]any{
										"type":        "string",
										"description": "The first argument.",
									},
									"arg2": map[string]any{
										"type":        "number",
										"description": "The second argument.",
									},
								},
								Required: []string{"arg1"},
							},
							Command:  "doTest.sh",
							Approval: "true",
						},
					},
					McpServer: map[string]iMcp.Server{
						"command1": {
							Command: iMcp.Command{
								Name:      "docker",
								Arguments: []string{"run", "--rm", "-i", "-e", "GITHUB_PERSONAL_ACCESS_TOKEN=github_", "ghcr.io/github/github-mcp-server"},
								Environment: map[string]string{
									"TEST": "test",
								},
							},
							Timeout: iMcp.Timeout{
								Init:      yacl.P(500 * time.Millisecond),
								List:      yacl.P(10 * time.Second),
								Execution: yacl.P(1*time.Minute + 30*time.Second),
							},
						},
						"command2": {
							Command: iMcp.Command{
								Name: "echo",
							},
						},
						"http": {
							Http: iMcp.Http{
								BaseUrl: "http://localhost:8080/api/v1",
								Headers: map[string]string{
									"Authorization": "Bearer TOKEN",
								},
							},
						},
					},
				},
			},
			Printer: model.PrinterConfig{
				Format:     "json",
				TargetsRaw: []string{"stdout"},
			},
		},
		DebugConfig: model.DebugConfig{
			LogLevel:     model.LogLevelDebug,
			PprofAddress: ":1312",
			VueDevTools: model.VueDevToolsConfig{
				Host: "localhost",
				Port: 1312,
			},
			WebKit: model.WebKitInspectorConfig{
				OpenInspectorOnStartup: true,
				HttpServerAddress:      "127.0.0.1:5000",
			},
		},
		Themes: model.Themes{
			Dark: &model.Theme{
				Colors: map[string]string{
					"background": "#FFFFFF",
					"surface":    "#FFFFFF",
				},
			},
			Light: &model.Theme{
				Colors: map[string]string{
					"chat-tool-call": "#FF0000",
				},
			},
			Custom: map[string]model.Theme{
				"test": {
					Colors: map[string]string{
						"background":     "#00FF00",
						"chat-tool-call": "#0000FF",
					},
				},
			},
		},
	}

	expConfig := expDefConf
	expConfig.Profiles = map[string]*model.Profile{
		"test": &expDefConf.MainProfile,
	}

	assert.Equal(t, &expConfig, c)
}
