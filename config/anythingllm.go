package config

import (
	"fmt"
	flag "github.com/spf13/pflag"
)

type AnythingLLMConfig struct {
	BaseURL   string
	Token     string
	Workspace string
}

func configureAnythingLLM(c *AnythingLLMConfig) {
	flag.StringVar(&c.BaseURL, "anythingllm-base-url", "", "Base URL for AnythingLLM")
	flag.StringVar(&c.Token, "anythingllm-token", "", "Token for AnythingLLM")
	flag.StringVar(&c.Workspace, "anythingllm-workspace", "", "Workspace for AnythingLLM")
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
