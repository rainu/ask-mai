package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type AnthropicConfig struct {
	Token   common.Secret `yaml:"api-key" usage:"API Key"`
	BaseUrl string        `yaml:"base-url" usage:"BaseUrl"`
	Model   string        `yaml:"model" usage:"Model"`
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

func (c *AnthropicConfig) Validate() error {
	if ce := c.Token.Validate(); ce != nil {
		return fmt.Errorf("Anthropic API Key is missing: %w", ce)
	}
	if c.Model == "" {
		return fmt.Errorf("Anthropic Model is missing")
	}

	return nil
}

func (c *AnthropicConfig) BuildLLM() (llms.Model, error) {
	return llms.NewAnthropic(c.AsOptions())
}
