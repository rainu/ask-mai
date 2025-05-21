package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/tools"
)

type CommandExecutionArguments struct {
	Name             string            `json:"name"`
	Arguments        []string          `json:"arguments"`
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`
}

var CommandExecutionDefinition = tools.BuiltinDefinition{
	Description: "Execute a command on the user's system.",
	Parameter: map[string]any{
		"type":   "object",
		"strict": true,
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "The name / path to the command to execute.",
			},
			"arguments": map[string]any{
				"type":        "array",
				"description": "The arguments for the command.",
				"items": map[string]any{
					"type": "string",
				},
			},
			"working_directory": map[string]any{
				"type":        "string",
				"description": "The working directory for the command.",
			},
			"environment": map[string]any{
				"type":                 "object",
				"description":          "Additional environment variables to pass to the command.",
				"additionalProperties": true,
			},
		},
		"required": []string{"name"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs CommandExecutionArguments
		err := json.Unmarshal([]byte(jsonArguments), &pArgs)
		if err != nil {
			return nil, fmt.Errorf("error parsing arguments: %w", err)
		}

		if pArgs.Name == "" {
			return nil, fmt.Errorf("missing parameter: 'name'")
		}

		cmdDesc := CommandDescriptor{
			Command:               pArgs.Name,
			Arguments:             pArgs.Arguments,
			AdditionalEnvironment: pArgs.Environment,
			WorkingDirectory:      pArgs.WorkingDirectory,
		}

		return cmdDesc.Run(ctx)
	},
}
