package llm

import (
	"fmt"
)

type AnythingLLMConfig struct {
	BaseURL   string `config:"base-url" usage:"Base URL"`
	Token     string `config:"token" usage:"Token"`
	Workspace string `config:"workspace" usage:"Workspace"`
}

func (c *AnythingLLMConfig) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("AnythingLLM Base URL is missing")
	}
	if c.Token == "" {
		return fmt.Errorf("AnythingLLM Token is missing")
	}
	if c.Workspace == "" {
		return fmt.Errorf("AnythingLLM Workspace is missing")
	}

	return nil
}
