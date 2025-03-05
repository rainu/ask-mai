package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cmdchain "github.com/rainu/go-command-chain"
)

type CommandExecution struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `yaml:"no-approval" json:"no-approval" usage:"Needs no user approval to be executed"`
}

func (c CommandExecution) AsFunctionDefinition() *FunctionDefinition {
	if c.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "executeCommand",
		Description: "Execute a command on the user's system.",
		Parameters: map[string]any{
			"type": "object",
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
		CommandFn:     c.Command,
		NeedsApproval: !c.NoApproval,
	}
}

type commandArguments struct {
	Name             string            `json:"name"`
	Arguments        []string          `json:"arguments"`
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`
}

func (c CommandExecution) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs commandArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if pArgs.Name == "" {
		return nil, fmt.Errorf("missing parameter: 'name'")
	}

	cmd := cmdchain.Builder().JoinWithContext(ctx, pArgs.Name, pArgs.Arguments...)

	if pArgs.WorkingDirectory != "" {
		cmd = cmd.WithWorkingDirectory(pArgs.WorkingDirectory)
	}
	if len(pArgs.Environment) > 0 {
		envMap := map[any]any{}
		for k, v := range pArgs.Environment {
			envMap[k] = v
		}
		cmd = cmd.WithAdditionalEnvironmentMap(envMap)
	}

	buf := bytes.NewBuffer([]byte{})
	execErr := cmd.Finalize().
		WithOutput(buf).
		WithError(buf).
		Run()

	return buf.Bytes(), execErr
}
