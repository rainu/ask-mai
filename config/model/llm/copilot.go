package llm

import (
	"fmt"
	"github.com/rainu/ask-mai/llms"
)

type CopilotConfig struct {
}

func (c *CopilotConfig) Validate() error {
	if !llms.IsCopilotInstalled() {
		return fmt.Errorf("GitHub Copilot is not installed")
	}
	return nil
}

func (c *CopilotConfig) BuildLLM() (llms.Model, error) {
	return llms.NewCopilot()
}
