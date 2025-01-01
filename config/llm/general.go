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

type LLMConfig struct {
	Backend string `config:"backend" short:"b"`

	LocalAI     LocalAIConfig     `config:"localai" usage:"LocalAI: "`
	OpenAI      OpenAIConfig      `config:"openai" usage:"OpenAI: "`
	AnythingLLM AnythingLLMConfig `config:"anythingllm" usage:"AnythingLLM: "`
	Ollama      OllamaConfig      `config:"ollama" usage:"Ollama: "`
	Mistral     MistralConfig     `config:"mistral" usage:"Mistral: "`
	Anthropic   AnthropicConfig   `config:"anthropic" usage:"Anthropic: "`
	CallOptions CallOptionsConfig `config:"call" usage:"LLM-CALL: "`
}

func (c *LLMConfig) GetUsage(field string) string {
	switch field {
	case "Backend":
		return fmt.Sprintf("The backend to use ('%s', '%s', '%s', '%s', '%s', '%s', '%s')", BackendCopilot, BackendOpenAI, BackendLocalAI, BackendAnythingLLM, BackendOllama, BackendMistral, BackendAnthropic)
	}
	return ""
}

func (c *LLMConfig) Validate() error {
	switch c.Backend {
	case BackendCopilot:
		if !llms.IsCopilotInstalled() {
			return fmt.Errorf("GitHub Copilot is not installed")
		}
	case BackendLocalAI:
		if ve := c.LocalAI.Validate(); ve != nil {
			return ve
		}
	case BackendOpenAI:
		if ve := c.OpenAI.Validate(); ve != nil {
			return ve
		}
	case BackendAnythingLLM:
		if ve := c.AnythingLLM.Validate(); ve != nil {
			return ve
		}
	case BackendOllama:
		if ve := c.Ollama.Validate(); ve != nil {
			return ve
		}
	case BackendMistral:
		if ve := c.Mistral.Validate(); ve != nil {
			return ve
		}
	case BackendAnthropic:
		if ve := c.Anthropic.Validate(); ve != nil {
			return ve
		}
	default:
		return fmt.Errorf("Invalid backend")
	}

	if ve := c.CallOptions.Validate(); ve != nil {
		return ve
	}

	return nil
}

func (c *LLMConfig) BuildLLM() (llms.Model, error) {
	switch c.Backend {
	case BackendCopilot:
		return llms.NewCopilot()
	case BackendLocalAI:
		return llms.NewLocalAI(
			c.LocalAI.AsOptions(),
		)
	case BackendOpenAI:
		return llms.NewOpenAI(
			c.OpenAI.AsOptions(),
		)
	case BackendAnythingLLM:
		return llms.NewAnythingLLM(
			c.AnythingLLM.BaseURL,
			c.AnythingLLM.Token,
			c.AnythingLLM.Workspace,
		)
	case BackendOllama:
		return llms.NewOllama(
			c.Ollama.AsOptions(),
		)
	case BackendMistral:
		return llms.NewMistral(
			c.Mistral.AsOptions(),
		)
	case BackendAnthropic:
		return llms.NewAnthropic(
			c.Anthropic.AsOptions(),
		)
	default:
		return nil, fmt.Errorf("unknown backend: %s", c.Backend)
	}
}
