package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/common"
	llmCommon "github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/ask-mai/internal/llms/localai"
	"github.com/tmc/langchaingo/llms/openai"
)

type LocalAIConfig struct {
	APIKey  common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	Model   string        `yaml:"model,omitempty" usage:"Model"`
	BaseUrl string        `yaml:"base-url,omitempty" usage:"BaseUrl"`
}

func (c *LocalAIConfig) AsOptions() (opts []openai.Option) {
	if v := c.APIKey.GetOrPanicWithDefaultTimeout(); v != nil && len(v) > 0 {
		opts = append(opts, openai.WithToken(string(v)))
	} else {
		// the underlying openai implementation wants to have an API key
		// so we'll just use a placeholder here
		// LocalAI doesn't need an API key and don't care about it
		opts = append(opts, openai.WithToken("PLACEHOLDER"))
	}
	if c.Model != "" {
		opts = append(opts, openai.WithModel(c.Model))
	}
	if c.BaseUrl != "" {
		opts = append(opts, openai.WithBaseURL(c.BaseUrl))
	}

	return
}

func (c *LocalAIConfig) Validate() error {
	if c.BaseUrl == "" {
		return fmt.Errorf("LocalAI Base URL is missing")
	}
	if c.Model == "" {
		return fmt.Errorf("LocalAI Model is missing")
	}

	return nil
}

func (c *LocalAIConfig) BuildLLM() (llmCommon.Model, error) {
	return localai.New(c.AsOptions())
}
