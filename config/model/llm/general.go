package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/llm/tools"
	"github.com/rainu/ask-mai/llms"
	"reflect"
	"slices"
	"strings"
)

type llmConfig interface {
	BuildLLM() (llms.Model, error)
	Validate() error
}

type LLMConfig struct {
	Backend string `yaml:"backend" short:"b"`

	Copilot     CopilotConfig     `yaml:"copilot" usage:"Copilot: " llm:""`
	LocalAI     LocalAIConfig     `yaml:"localai" usage:"LocalAI: " llm:""`
	OpenAI      OpenAIConfig      `yaml:"openai" usage:"OpenAI: " llm:""`
	AnythingLLM AnythingLLMConfig `yaml:"anythingllm" usage:"AnythingLLM: " llm:""`
	Ollama      OllamaConfig      `yaml:"ollama" usage:"Ollama: " llm:""`
	Mistral     MistralConfig     `yaml:"mistral" usage:"Mistral: " llm:""`
	Anthropic   AnthropicConfig   `yaml:"anthropic" usage:"Anthropic: " llm:""`
	DeepSeek    DeepSeekConfig    `yaml:"deepseek" usage:"DeepSeek: " llm:""`

	CallOptions CallOptionsConfig `yaml:"call" usage:"LLM-CALL: "`
	Tools       tools.Config      `yaml:"tool" usage:"LLM-TOOLS: "`
}

func (c *LLMConfig) getBackend() llmConfig {
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := strings.ToLower(val.Type().Field(i).Name)
		if name == strings.ToLower(c.Backend) {
			return val.Field(i).Addr().Interface().(llmConfig)
		}
	}
	return nil
}

func (c *LLMConfig) listBackends() (result []string) {
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		_, ok := f.Tag.Lookup("llm")
		if ok {
			result = append(result, strings.ToLower(f.Name))
		}
	}
	slices.Sort(result)
	return
}

func (c *LLMConfig) GetUsage(field string) string {
	switch field {
	case "Backend":
		return fmt.Sprintf("The backend to use (%s)", strings.Join(c.listBackends(), ", "))
	}
	return ""
}

func (c *LLMConfig) Validate() error {
	b := c.getBackend()
	if b == nil {
		return fmt.Errorf("Invalid backend '%s'", c.Backend)
	}
	if ve := b.Validate(); ve != nil {
		return ve
	}
	if ve := c.CallOptions.Validate(); ve != nil {
		return ve
	}
	if ve := c.Tools.Validate(); ve != nil {
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
