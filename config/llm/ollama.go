package llm

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaConfig struct {
	ServerURL string `config:"server-url" usage:"Server URL"`
	Model     string `config:"model" usage:"Model"`
}

func configureOllama(c *OllamaConfig) {
	flag.StringVar(&c.ServerURL, "ollama-server-url", "", "Server URL for Ollama")
	flag.StringVar(&c.Model, "ollama-model", "", "Model for Ollama")
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
