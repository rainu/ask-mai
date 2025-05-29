package llm

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/ask-mai/internal/llms/copilot"
	langLLMS "github.com/tmc/langchaingo/llms"
	"reflect"
	"slices"
	"strings"
)

type llmConfig interface {
	BuildLLM() (common.Model, error)
	Validate() error
}

type LLMConfig struct {
	Backend string `yaml:"backend,omitempty" short:"b"`

	Copilot     CopilotConfig     `yaml:"copilot,omitempty" usage:"Copilot " llm:""`
	LocalAI     LocalAIConfig     `yaml:"localai,omitempty" usage:"LocalAI " llm:""`
	OpenAI      OpenAIConfig      `yaml:"openai,omitempty" usage:"OpenAI " llm:""`
	AnythingLLM AnythingLLMConfig `yaml:"anythingllm,omitempty" usage:"AnythingLLM " llm:""`
	Ollama      OllamaConfig      `yaml:"ollama,omitempty" usage:"Ollama " llm:""`
	Mistral     MistralConfig     `yaml:"mistral,omitempty" usage:"Mistral " llm:""`
	Anthropic   AnthropicConfig   `yaml:"anthropic,omitempty" usage:"Anthropic " llm:""`
	DeepSeek    DeepSeekConfig    `yaml:"deepseek,omitempty" usage:"DeepSeek " llm:""`

	CallOptions CallOptionsConfig `yaml:"call,omitempty" usage:"Call option "`
	Tool        tools.Config      `yaml:"tool,omitempty"`
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

func (c *LLMConfig) SetDefaults() {
	if c.Backend == "" {
		if copilot.IsCopilotInstalled() {
			c.Backend = "copilot"
		}
	}
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
	if ve := c.Tool.Validate(); ve != nil {
		return ve
	}

	return nil
}

func (c *LLMConfig) BuildLLM() (common.Model, error) {
	b := c.getBackend()
	if b == nil {
		return nil, fmt.Errorf("unknown backend: %s", c.Backend)
	}
	return b.BuildLLM()
}

func (c *LLMConfig) AsOptions(ctx context.Context) ([]langLLMS.CallOption, error) {
	mcpTools, err := c.Tool.GetTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list mcp tools: %w", err)
	}

	var tools []langLLMS.Tool
	for key, toolDef := range mcpTools {
		tools = append(tools, langLLMS.Tool{
			Type: "function",
			Function: &langLLMS.FunctionDefinition{
				Name:        key,
				Description: toolDef.Description,
				Parameters:  toolDef.InputSchema,
			},
		})
	}

	// sort tools by name so that the order is consistent
	// and the llm can cache the tool definitions
	slices.SortFunc(tools, func(a, b langLLMS.Tool) int {
		return strings.Compare(a.Function.Name, b.Function.Name)
	})

	opts := c.CallOptions.AsOptions()
	if len(tools) > 0 {
		opts = append(opts, langLLMS.WithTools(tools))
	}
	return opts, nil
}
