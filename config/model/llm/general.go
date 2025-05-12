package llm

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/config/model/llm/mcp"
	"github.com/rainu/ask-mai/config/model/llm/tools"
	"github.com/rainu/ask-mai/llms"
	langLLMS "github.com/tmc/langchaingo/llms"
	"reflect"
	"slices"
	"strings"
)

type llmConfig interface {
	BuildLLM() (llms.Model, error)
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

	CallOptions CallOptionsConfig     `yaml:"call,omitempty" usage:"Call option "`
	Tools       tools.Config          `yaml:"tool,omitempty"`
	McpServer   map[string]mcp.Server `yaml:"mcpServers,omitempty" usage:"MCP server "`
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
		if llms.IsCopilotInstalled() {
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
	if ve := c.Tools.Validate(); ve != nil {
		return ve
	}
	for name, server := range c.McpServer {
		if ve := server.Validate(); ve != nil {
			return fmt.Errorf("invlalid mcpServer config for %s: %w", name, ve)
		}
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

func (c *LLMConfig) AsOptions(ctx context.Context) ([]langLLMS.CallOption, error) {
	var tools []langLLMS.Tool

	for name, definition := range c.Tools.GetTools() {
		tools = append(tools, langLLMS.Tool{
			Type: "function",
			Function: &langLLMS.FunctionDefinition{
				Name:        name,
				Description: definition.Description,
				Parameters:  definition.Parameters,
			},
		})
	}

	mcpTools, err := mcp.MergeTools(ctx, c.McpServer)
	if err != nil {
		return nil, fmt.Errorf("failed to list mcp tools: %w", err)
	}

	for key, toolDef := range mcpTools {
		desc := ""
		if toolDef.Description != nil {
			desc = *toolDef.Description
		}

		tools = append(tools, langLLMS.Tool{
			Type: "function",
			Function: &langLLMS.FunctionDefinition{
				Name:        key,
				Description: desc,
				Parameters:  toolDef.InputSchema,
			},
		})
	}

	opts := c.CallOptions.AsOptions()
	if len(tools) > 0 {
		opts = append(opts, langLLMS.WithTools(tools))
	}
	return opts, nil
}
