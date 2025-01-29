package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type DeepSeekConfig struct {
	APIKey  string `yaml:"api-key" usage:"API Key"`
	Model   string `yaml:"model" usage:"Model"`
	BaseUrl string `yaml:"base-url" usage:"BaseUrl"`
}

func (c *DeepSeekConfig) AsOptions() (opts []openai.Option) {
	if c.APIKey != "" {
		opts = append(opts, openai.WithToken(c.APIKey))
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

func (c *DeepSeekConfig) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("DeepSeek API Key is missing")
	}
	if c.Model == "" {
		return fmt.Errorf("DeepSeek Model is missing")
	}

	return nil
}

func (c *DeepSeekConfig) BuildLLM() (llms.Model, error) {
	return llms.NewDeepSeek(c.AsOptions())
}
