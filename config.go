package main

import (
	"flag"
	"fmt"
	"github.com/rainu/ask-mai/backend/copilot"
)

const (
	BackendCopilot     = "copilot"
	BackendOpenAI      = "openai"
	BackendAnythingLLM = "anythingllm"
)

type Config struct {
	Prompt string

	Backend     string
	OpenAI      OpenAIConfig
	AnythingLLM AnythingLLMConfig
}

type OpenAIConfig struct {
	APIKey string
}

type AnythingLLMConfig struct {
	BaseURL   string
	Token     string
	Workspace string
}

func ParseConfig() Config {
	var config Config

	flag.StringVar(&config.Prompt, "prompt", "", "The prompt to use")
	flag.StringVar(&config.Backend, "backend", BackendCopilot, fmt.Sprintf("The backend to use ('%s', '%s', '%s')", BackendCopilot, BackendOpenAI, BackendAnythingLLM))
	flag.StringVar(&config.OpenAI.APIKey, "openai-api-key", "", "OpenAI API Key")
	flag.StringVar(&config.AnythingLLM.BaseURL, "anythingllm-base-url", "", "Base URL for AnythingLLM")
	flag.StringVar(&config.AnythingLLM.Token, "anythingllm-token", "", "Token for AnythingLLM")
	flag.StringVar(&config.AnythingLLM.Workspace, "anythingllm-workspace", "", "Workspace for AnythingLLM")

	flag.Parse()

	return config
}

func (c Config) Validate() error {
	if c.Backend != BackendCopilot && c.Backend != BackendOpenAI && c.Backend != BackendAnythingLLM {
		return fmt.Errorf("Invalid backend")
	}

	if c.Backend == BackendCopilot && !copilot.IsCopilotInstalled() {
		return fmt.Errorf("GitHub Copilot is not installed")
	}

	if c.Backend == BackendOpenAI && c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API Key is missing")
	}

	if c.Backend == BackendAnythingLLM {
		if c.AnythingLLM.BaseURL == "" {
			return fmt.Errorf("AnythingLLM Base URL is missing")
		}
		if c.AnythingLLM.Token == "" {
			return fmt.Errorf("AnythingLLM Token is missing")
		}
		if c.AnythingLLM.Workspace == "" {
			return fmt.Errorf("AnythingLLM Workspace is missing")
		}
	}

	return nil
}
