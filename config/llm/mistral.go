package llm

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/tmc/langchaingo/llms/mistral"
)

type MistralConfig struct {
	ApiKey   string `config:"api-key" usage:"API Key"`
	Endpoint string `config:"endpoint" usage:"Endpoint"`
	Model    string `config:"model" usage:"Model"`
}

func configureMistral(c *MistralConfig) {
	flag.StringVar(&c.ApiKey, "mistral-api-key", "", "API Key for Mistral")
	flag.StringVar(&c.Endpoint, "mistral-endpoint", "", "Endpoint for Mistral")
	flag.StringVar(&c.Model, "mistral-model", "", "Model for Mistral")
}

func (c *MistralConfig) AsOptions() (opts []mistral.Option) {
	if c.ApiKey != "" {
		opts = append(opts, mistral.WithAPIKey(c.ApiKey))
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
	if c.ApiKey == "" {
		return fmt.Errorf("Mistral API Key is missing")
	}

	return nil
}
