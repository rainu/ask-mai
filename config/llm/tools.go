package llm

import (
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
	"mvdan.cc/sh/v3/shell"
)

const FunctionArgumentNameAll = "@"

type ToolsConfig struct {
	RawTools []string `config:"function" yaml:"-" usage:"Function definition (json) to use. See Tool-Help (--help-tool) for more information."`

	Tools map[string]FunctionDefinition `config:"-" yaml:"functions"`
}

type FunctionDefinition struct {
	Name          string `config:"name" yaml:"-" json:"name" usage:"The name of the function"`
	Description   string `yaml:"description" json:"description" usage:"The description of the function"`
	Parameters    any    `yaml:"parameters" json:"parameters" usage:"The parameter definition of the function"`
	Command       string `yaml:"command" json:"command" usage:"The command to execute. This is a format string with placeholders for the parameters. Example: /usr/bin/touch $path"`
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

func (f *FunctionDefinition) GetCommandWithArgs(jsonArgs string) (string, []string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonArgs), &data); err != nil {
		return "", nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fields, err := shell.Fields(f.Command, func(varName string) string {
		if varName == FunctionArgumentNameAll {
			return jsonArgs
		}

		varValue, exists := data[varName]
		if !exists {
			return ""
		}

		val, err := json.Marshal(varValue)
		if err != nil {
			slog.Error("Failed to marshal value",
				"varName", varName,
				"value", data[varName],
				"error", err,
			)
			return ""
		}
		sVal := string(val)
		if len(sVal) > 0 && sVal[0] == '"' {
			sVal = sVal[1:]
		}
		if len(sVal) > 0 && sVal[len(sVal)-1] == '"' {
			sVal = sVal[:len(sVal)-1]
		}
		return sVal
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse command: %w", err)
	}
	return fields[0], fields[1:], nil
}
