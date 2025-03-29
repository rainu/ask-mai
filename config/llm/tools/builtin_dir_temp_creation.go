package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type DirectoryTempCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y DirectoryTempCreationResult    `config:"-" yaml:"-"`
	Z DirectoryTempCreationArguments `config:"-" yaml:"-"`
}

func (f DirectoryTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createTempDirectory",
		Description: "Creates a new temporary directory on the user's system.",
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

type DirectoryTempCreationArguments struct {
}

type DirectoryTempCreationResult struct {
	Path string `json:"path"`
}

func (f DirectoryTempCreation) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs DirectoryTempCreationArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	path, err := os.MkdirTemp("", "ask-mai.*")
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	return json.Marshal(DirectoryTempCreationResult{
		Path: path,
	})
}
