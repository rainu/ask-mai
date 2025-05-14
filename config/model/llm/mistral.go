package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	llmCommon "github.com/rainu/ask-mai/llms/common"
	iMistral "github.com/rainu/ask-mai/llms/mistral"
	"github.com/tmc/langchaingo/llms/mistral"
)

type MistralConfig struct {
	ApiKey   common.Secret `yaml:"api-key,omitempty" usage:"API Key"`
	Endpoint string        `yaml:"endpoint,omitempty" usage:"Endpoint"`
	Model    string        `yaml:"model,omitempty" usage:"Model"`
}

func (c *MistralConfig) AsOptions() (opts []mistral.Option) {
	if v := c.ApiKey.GetOrPanicWithDefaultTimeout(); v != nil {
		opts = append(opts, mistral.WithAPIKey(string(v)))
	}
	if c.Endpoint != "" {
		opts = append(opts, mistral.WithEndpoint(c.Endpoint))
	}
	if c.Model != "" {
		opts = append(opts, mistral.WithModel(c.Model))
	}

	return
}

func (c *MistralConfig) Validate() error {
	if ce := c.ApiKey.Validate(); ce != nil {
		return fmt.Errorf("Mistral API Key is missing: %w", ce)
	}

	return nil
}

func (c *MistralConfig) BuildLLM() (llmCommon.Model, error) {
	return iMistral.New(c.AsOptions())
}
