package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms/common"
	iOllama "github.com/rainu/ask-mai/llms/ollama"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaConfig struct {
	ServerURL string `yaml:"server-url,omitempty" usage:"Server URL"`
	Model     string `yaml:"model,omitempty" usage:"Model"`
}

func (c *OllamaConfig) AsOptions() (opts []ollama.Option) {
	if c.ServerURL != "" {
		opts = append(opts, ollama.WithServerURL(c.ServerURL))
	}
	if c.Model != "" {
		opts = append(opts, ollama.WithModel(c.Model))
	}

	return
}

func (c *OllamaConfig) Validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("Ollama Server URL is missing")
	}
	if c.Model == "" {
		return fmt.Errorf("Ollama Model is missing")
	}

	return nil
}

func (c *OllamaConfig) BuildLLM() (common.Model, error) {
	return iOllama.New(c.AsOptions())
}
