package llm

import (
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type CallOptionsConfig struct {
	SystemPrompt string  `yaml:"system-prompt" short:"S" usage:"System Prompt"`
	MaxToken     int     `yaml:"max-token" usage:"Max Token"`
	Temperature  float64 `yaml:"temperature" usage:"Temperature"`
	TopK         int     `yaml:"top-k" usage:"Top-K"`
	TopP         float64 `yaml:"top-p" usage:"Top-P"`
	MinLength    int     `yaml:"min-length" usage:"Min Length"`
	MaxLength    int     `yaml:"max-length" usage:"Max Length"`
}

func (c *CallOptionsConfig) AsOptions() (opts []llms.CallOption) {
	if c.MaxToken != 0 {
		opts = append(opts, llms.WithMaxTokens(c.MaxToken))
	}
	if c.Temperature != -1 {
		opts = append(opts, llms.WithTemperature(c.Temperature))
	}
	if c.TopK != -1 {
		opts = append(opts, llms.WithTopK(c.TopK))
	}
	if c.TopP != -1 {
		opts = append(opts, llms.WithTopP(c.TopP))
	}
	if c.MinLength != 0 {
		opts = append(opts, llms.WithMinLength(c.MinLength))
	}
	if c.MaxLength != 0 {
		opts = append(opts, llms.WithMaxLength(c.MaxLength))
	}

	// ask-mai can only handle one choice - to here we will "force" the llm to only generate one choice
	opts = append(opts, llms.WithN(1))

	return
}

func (c *CallOptionsConfig) Validate() error {
	if c.Temperature != -1 && (c.Temperature < 0 || c.Temperature > 1) {
		return fmt.Errorf("LLM-Call Temperature is invalid")
	}
	if c.TopK != -1 && c.TopK < 0 {
		return fmt.Errorf("LLM-Call Top-K is invalid")
	}
	if c.TopP != -1 && (c.TopP < 0 || c.TopP > 1) {
		return fmt.Errorf("LLM-Call Top-P is invalid")
	}
	if c.MinLength < 0 {
		return fmt.Errorf("LLM-Call Min Length is invalid")
	}
	if c.MaxLength < 0 {
		return fmt.Errorf("LLM-Call Max Length")
	}

	return nil
}
