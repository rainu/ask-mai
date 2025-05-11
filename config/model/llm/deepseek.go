package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type DeepSeekConfig struct {
	APIKey  common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	Model   string        `yaml:"model,omitempty" usage:"Model"`
	BaseUrl string        `yaml:"base-url,omitempty" usage:"BaseUrl"`
}

func (c *DeepSeekConfig) AsOptions() (opts []openai.Option) {
	if v := c.APIKey.GetOrPanicWithDefaultTimeout(); v != nil {
		opts = append(opts, openai.WithToken(string(v)))
	}
	if c.Model != "" {
		opts = append(opts, openai.WithModel(c.Model))
	}
	if c.BaseUrl == "" {
		opts = append(opts, openai.WithBaseURL("https://api.deepseek.com/v1"))
	} else {
		opts = append(opts, openai.WithBaseURL(c.BaseUrl))
	}

	return
}

func (c *DeepSeekConfig) SetDefaults() {
	if c.Model == "" {
		c.Model = "deepseek-chat"
	}
}

func (c *DeepSeekConfig) Validate() error {
	if ce := c.APIKey.Validate(); ce != nil {
		return fmt.Errorf("DeepSeek API Key is missing: %w", ce)
	}
	if c.Model == "" {
		return fmt.Errorf("DeepSeek Model is missing")
	}

	return nil
}

func (c *DeepSeekConfig) BuildLLM() (llms.Model, error) {
	return llms.NewDeepSeek(c.AsOptions())
}
