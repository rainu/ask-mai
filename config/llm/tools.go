package llm

import (
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type ToolsConfig struct {
	RawTools []string `config:"function" yaml:"-" usage:"Function definition (json) to use. See Tool-Help (--help-tool) for more information."`

	Tools map[string]FunctionDefinition `config:"-" yaml:"functions"`
}

type FunctionDefinition struct {
	Name          string `config:"name" yaml:"-" json:"name" usage:"The name of the function"`
	Description   string `yaml:"description" json:"description" usage:"The description of the function"`
	Parameters    any    `yaml:"parameters" json:"parameters" usage:"The parameter definition of the function"`
	Command       string `yaml:"command" json:"command" usage:"The command to execute"`
	NeedsApproval bool   `yaml:"approval" json:"approval" usage:"Needs user approval to be executed"`
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
