package config

import (
	"flag"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type CallOptionsConfig struct {
	SystemPrompt string
	MaxToken     int
	Temperature  float64
	TopK         int
	TopP         float64
	MinLength    int
	MaxLength    int
}

func configureCallOptions(c *CallOptionsConfig) {
	flag.StringVar(&c.SystemPrompt, "call-system-prompt", "", "LLM-Call System Prompt")
	flag.IntVar(&c.MaxToken, "call-max-token", 0, "LLM-Call Max Token")
	flag.Float64Var(&c.Temperature, "call-temperature", -1, "LLM-Call Temperature")
	flag.IntVar(&c.TopK, "call-top-k", -1, "LLM-Call Top-K")
	flag.Float64Var(&c.TopP, "call-top-p", -1, "LLM-Call Top-P")
	flag.IntVar(&c.MinLength, "call-min-length", 0, "LLM-Call Min Length")
	flag.IntVar(&c.MaxLength, "call-max-length", 0, "LLM-Call Max Length")
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
