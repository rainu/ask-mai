package llm

import (
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type ToolsConfig struct {
	RawTools []string `config:"function" yaml:"-" usage:"Function definition (json) to use"`

	Tools map[string]FunctionDefinition `config:"-" yaml:"functions"`
}

type FunctionDefinition struct {
	Name          string `yaml:"-" json:"name"`
	Description   string `yaml:"description" json:"description"`
	Parameters    any    `yaml:"parameters" json:"parameters"`
	Command       string `yaml:"command" json:"command"`
	NeedsApproval bool   `yaml:"approval" json:"approval"`
}

func (t *ToolsConfig) Validate() error {
	for i, tool := range t.RawTools {
		var result FunctionDefinition

		err := json.Unmarshal([]byte(tool), &result)
		if err != nil {
			return fmt.Errorf("Invalid tool definition #%d: %w", i, err)
		}

		t.Tools[result.Name] = result
	}

	for cmd, definition := range t.Tools {
		if definition.Command == "" {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}
	}

	return nil
}

func (t *ToolsConfig) AsOptions() (opts []llms.CallOption) {
	var tools []llms.Tool
	for name, tool := range t.Tools {
		tools = append(tools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}

	if len(tools) > 0 {
		opts = append(opts, llms.WithTools(tools))
	}
	return
}
