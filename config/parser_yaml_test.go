package config

import (
	"github.com/rainu/ask-mai/config/expression"
	"github.com/rainu/ask-mai/config/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_processYaml(t *testing.T) {
	c := &Config{}

	sr := strings.NewReader(`
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
        api-key: APIKey
        model: model
        base-url: baseurl
    openai:
        api-key: APIKey
        api-type: APIType
        api-version: APIVersion
        model: Model
        base-url: BaseUrl
        organization: Organization
    anythingllm:
        base-url: BaseURL
        token: Token
        workspace: Workspace
    ollama:
        server-url: ServerURL
        model: Model
    mistral:
        api-key: ApiKey
        endpoint: Endpoint
        model: Model
    anthropic:
        api-key: Token
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
print:
    format: json
    targets: 
        - stdout
debug:
    log-level: 1
    pprof-address: ":1312"
    vue-dev-tools:
        host: "localhost"
        port: 1312
    webkit:
        open-inspector: true
        http-server: "127.0.0.1:5000"
`)

	require.NoError(t, processYaml(sr, c))
	assert.Equal(t, &Config{
		UI: UIConfig{
			Window: WindowConfig{
				Title:            "Test Window",
				InitialWidth:     expression.NumberContainer{Expression: "800"},
				MaxHeight:        expression.NumberContainer{Expression: "600"},
				InitialPositionX: expression.NumberContainer{Expression: "100"},
				InitialPositionY: expression.NumberContainer{Expression: "100"},
				InitialZoom:      expression.NumberContainer{Expression: "1.0"},
				BackgroundColor:  WindowBackgroundColor{R: 255, G: 255, B: 255, A: 255},
				StartState:       1,
				AlwaysOnTop:      true,
				Frameless:        true,
				Resizeable:       true,
				Translucent:      "never",
			},
			Prompt: PromptConfig{
				InitValue:       "Initial Prompt",
				InitAttachments: []string{"attachment1", "attachment2"},
				MinRows:         1,
				MaxRows:         10,
				SubmitShortcut: Shortcut{
					Binding: []string{"alt+ctrl+meta+shift+enter"},
				},
			},
			FileDialog: FileDialogConfig{
				DefaultDirectory:           "/root",
				ShowHiddenFiles:            true,
				CanCreateDirectories:       true,
				ResolveAliases:             true,
				TreatPackagesAsDirectories: true,
				FilterDisplay:              []string{"Image"},
				FilterPattern:              []string{"*.png"},
			},
			Stream: true,
			QuitShortcut: Shortcut{
				Binding: []string{"alt+ctrl+meta+shift+escape"},
			},
			Theme:     "dark",
			CodeStyle: "default",
			Language:  "en",
		},
		LLM: llm.LLMConfig{
			Backend: "anthropic",
			LocalAI: llm.LocalAIConfig{
				APIKey:  "APIKey",
				Model:   "model",
				BaseUrl: "baseurl",
			},
			OpenAI: llm.OpenAIConfig{
				APIKey:       "APIKey",
				APIType:      "APIType",
				APIVersion:   "APIVersion",
				Model:        "Model",
				BaseUrl:      "BaseUrl",
				Organization: "Organization",
			},
			AnythingLLM: llm.AnythingLLMConfig{
				BaseURL:   "BaseURL",
				Token:     "Token",
				Workspace: "Workspace",
			},
			Ollama: llm.OllamaConfig{
				ServerURL: "ServerURL",
				Model:     "Model",
			},
			Mistral: llm.MistralConfig{
				ApiKey:   "ApiKey",
				Endpoint: "Endpoint",
				Model:    "Model",
			},
			Anthropic: llm.AnthropicConfig{
				Token:   "Token",
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
		},
		Printer: PrinterConfig{
			Format:     "json",
			TargetsRaw: []string{"stdout"},
		},
		Debug: DebugConfig{
			LogLevel:     1,
			PprofAddress: ":1312",
			VueDevTools: VueDevToolsConfig{
				Host: "localhost",
				Port: 1312,
			},
			WebKit: WebKitInspectorConfig{
				OpenInspectorOnStartup: true,
				HttpServerAddress:      "127.0.0.1:5000",
			},
		},
	}, c)
}
