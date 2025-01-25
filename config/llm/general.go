package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms"
)

const (
	BackendCopilot     = "copilot"
	BackendLocalAI     = "localai"
	BackendOpenAI      = "openai"
	BackendAnythingLLM = "anythingllm"
	BackendOllama      = "ollama"
	BackendMistral     = "mistral"
	BackendAnthropic   = "anthropic"
)

type llmConfig interface {
	BuildLLM() (llms.Model, error)
	Validate() error
}

type LLMConfig struct {
	Backend string `yaml:"backend" short:"b"`

	Copilot     CopilotConfig     `yaml:"copilot" usage:"Copilot: "`
	LocalAI     LocalAIConfig     `yaml:"localai" usage:"LocalAI: "`
	OpenAI      OpenAIConfig      `yaml:"openai" usage:"OpenAI: "`
	AnythingLLM AnythingLLMConfig `yaml:"anythingllm" usage:"AnythingLLM: "`
	Ollama      OllamaConfig      `yaml:"ollama" usage:"Ollama: "`
	Mistral     MistralConfig     `yaml:"mistral" usage:"Mistral: "`
	Anthropic   AnthropicConfig   `yaml:"anthropic" usage:"Anthropic: "`
	CallOptions CallOptionsConfig `yaml:"call" usage:"LLM-CALL: "`
}

func (c *LLMConfig) getBackend() llmConfig {
	switch c.Backend {
	case BackendCopilot:
		return &c.Copilot
	case BackendLocalAI:
		return &c.LocalAI
	case BackendOpenAI:
		return &c.OpenAI
	case BackendAnythingLLM:
		return &c.AnythingLLM
	case BackendOllama:
		return &c.Ollama
	case BackendMistral:
		return &c.Mistral
	case BackendAnthropic:
		return &c.Anthropic
	default:
		return nil
	}
}

func (c *LLMConfig) GetUsage(field string) string {
	switch field {
	case "Backend":
		return fmt.Sprintf("The backend to use ('%s', '%s', '%s', '%s', '%s', '%s', '%s')", BackendCopilot, BackendOpenAI, BackendLocalAI, BackendAnythingLLM, BackendOllama, BackendMistral, BackendAnthropic)
	}
	return ""
}

func (c *LLMConfig) Validate() error {
	b := c.getBackend()
	if b == nil {
		return fmt.Errorf("Invalid backend %s", c.Backend)
	}
	if ve := b.Validate(); ve != nil {
		return ve
	}
	if ve := c.CallOptions.Validate(); ve != nil {
		return ve
	}

	return nil
}

func (c *LLMConfig) BuildLLM() (llms.Model, error) {
	b := c.getBackend()
	if b == nil {
		return nil, fmt.Errorf("unknown backend: %s", c.Backend)
	}
	return b.BuildLLM()
}
