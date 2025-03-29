package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Environment struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y EnvironmentResult    `config:"-" yaml:"-"`
	Z EnvironmentArguments `config:"-" yaml:"-"`
}

func (f Environment) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getEnvironment",
		Description: "Get all environment variables of the user's system.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
			"type":                 "object",
			"properties":           map[string]any{},
			"additionalProperties": false,
			"required":             []string{},
		},
		NeedsApproval: f.NeedsApproval,
	}
}

type EnvironmentArguments struct {
}

type EnvironmentResult struct {
	Environment map[string]string `json:"env"`
}

func (f Environment) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs EnvironmentArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	result := EnvironmentResult{
		Environment: map[string]string{},
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		result.Environment[pair[0]] = pair[1]
	}

	return json.Marshal(result)
}
