package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model/common"
	"github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms/mistral"
)

type MistralConfig struct {
	ApiKey   common.Secret `yaml:"api-key" usage:"API Key"`
	Endpoint string        `yaml:"endpoint" usage:"Endpoint"`
	Model    string        `yaml:"model" usage:"Model"`
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

func (c *MistralConfig) BuildLLM() (llms.Model, error) {
	return llms.NewMistral(c.AsOptions())
}
