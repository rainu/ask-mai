package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms"
)

type AnythingLLMConfig struct {
	BaseURL      string `yaml:"base-url" usage:"Base URL"`
	Token        string `yaml:"token" usage:"Token"`
	Workspace    string `yaml:"workspace" usage:"Workspace"`
	DeleteThread bool   `yaml:"delete-thread" usage:"Delete the thread after the session is closed"`
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

func (c *AnythingLLMConfig) BuildLLM() (llms.Model, error) {
	return llms.NewAnythingLLM(
		c.BaseURL,
		c.Token,
		c.Workspace,
		c.DeleteThread,
	)
}
