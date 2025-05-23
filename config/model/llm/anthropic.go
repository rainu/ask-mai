package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	iAnthropic "github.com/rainu/ask-mai/llms/anthropic"
	llmCommon "github.com/rainu/ask-mai/llms/common"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type AnthropicConfig struct {
	Token   common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	BaseUrl string        `yaml:"base-url,omitempty" usage:"BaseUrl"`
	Model   string        `yaml:"model,omitempty" usage:"Model"`
}

func (c *AnthropicConfig) AsOptions() (opts []anthropic.Option) {
	if v := c.Token.GetOrPanicWithDefaultTimeout(); v != nil {
		opts = append(opts, anthropic.WithToken(string(v)))
	}
	if c.BaseUrl != "" {
		opts = append(opts, anthropic.WithBaseURL(c.BaseUrl))
	}
	if c.Model != "" {
		opts = append(opts, anthropic.WithModel(c.Model))
	}

	return
}

func (c *AnthropicConfig) SetDefaults() {
	if c.Model == "" {
		c.Model = "claude-3-5-haiku-latest"
	}
}

func (c *AnthropicConfig) Validate() error {
	if ce := c.Token.Validate(); ce != nil {
		return fmt.Errorf("Anthropic API Key is missing: %w", ce)
	}
	if c.Model == "" {
		return fmt.Errorf("Anthropic Model is missing")
	}

	return nil
}

func (c *AnthropicConfig) BuildLLM() (llmCommon.Model, error) {
	return iAnthropic.New(c.AsOptions())
}
