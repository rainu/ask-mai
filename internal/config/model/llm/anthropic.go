package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/common"
	iAnthropic "github.com/rainu/ask-mai/internal/llms/anthropic"
	llmCommon "github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/go-yacl"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type AnthropicConfig struct {
	Token   common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	BaseUrl string        `yaml:"base-url,omitempty" usage:"BaseUrl"`
	Model   string        `yaml:"model,omitempty" usage:"Model"`

	Cache AnthropicCache `yaml:"disable-cache,omitempty" usage:"disable "`
}

type AnthropicCache struct {
	SystemMessage *bool `yaml:"system-message,omitempty" usage:"system message cache"`
	Tools         *bool `yaml:"tools,omitempty" usage:"tools cache"`
	Chat          *bool `yaml:"chat,omitempty" usage:"chat cache"`
}

func (c *AnthropicCache) SetDefaults() {
	if c.SystemMessage == nil {
		c.SystemMessage = yacl.P(false)
	}
	if c.Tools == nil {
		c.Tools = yacl.P(false)
	}
	if c.Chat == nil {
		c.Chat = yacl.P(false)
	}
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
	if c.Cache.SystemMessage != nil && !*c.Cache.SystemMessage {
		opts = append(opts, anthropic.WithCacheSystemMessage())
	}
	if c.Cache.Tools != nil && !*c.Cache.Tools {
		opts = append(opts, anthropic.WithCacheTools())
	}
	if c.Cache.Chat != nil && !*c.Cache.Chat {
		opts = append(opts, anthropic.WithCacheChat())
	}

	return
}

func (c *AnthropicConfig) BuildLLM() (llmCommon.Model, error) {
	return iAnthropic.New(c.AsOptions())
}
