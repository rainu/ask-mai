package llm

import (
	"fmt"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type AnthropicConfig struct {
	Token   string `yaml:"api-key" usage:"API Key"`
	BaseUrl string `yaml:"base-url" usage:"BaseUrl"`
	Model   string `yaml:"model" usage:"Model"`
}

func (c *AnthropicConfig) AsOptions() (opts []anthropic.Option) {
	if c.Token != "" {
		opts = append(opts, anthropic.WithToken(c.Token))
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
	if c.Token == "" {
		return fmt.Errorf("Anthropic API Key is missing")
	}
	if c.Model == "" {
		return fmt.Errorf("Anthropic Model is missing")
	}

	return nil
}
