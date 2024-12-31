package config

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type AnthropicConfig struct {
	Token   string
	BaseUrl string
	Model   string
}

func configureAnthropic(c *AnthropicConfig) {
	flag.StringVar(&c.Token, "anthropic-api-key", "", "API Key for Anthropic")
	flag.StringVar(&c.BaseUrl, "anthropic-base-url", "", "BaseUrl for Anthropic")
	flag.StringVar(&c.Model, "anthropic-model", "", "Model for Anthropic")
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
