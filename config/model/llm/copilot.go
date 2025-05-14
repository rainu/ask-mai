package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms/common"
	"github.com/rainu/ask-mai/llms/copilot"
)

type CopilotConfig struct {
}

func (c *CopilotConfig) Validate() error {
	if !copilot.IsCopilotInstalled() {
		return fmt.Errorf("GitHub Copilot is not installed")
	}
	return nil
}

func (c *CopilotConfig) BuildLLM() (common.Model, error) {
	return copilot.New()
}
