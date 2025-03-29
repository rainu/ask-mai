package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type FileDeletion struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y FileDeletionResult    `config:"-" yaml:"-"`
	Z FileDeletionArguments `config:"-" yaml:"-"`
}

func (f FileDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "deleteFile",
		Description: "Delete a file on the user's system.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path to the file to delete. Use '~' as placeholder for the user's home directory.",
				},
			},
			"additionalProperties": false,
			"required":             []string{"path"},
		},
		NeedsApproval: !f.NoApproval,
	}
}

type FileDeletionArguments struct {
	Path Path `json:"path"`
}

type FileDeletionResult struct {
	Path string `json:"path"`
}

func (f FileDeletion) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs FileDeletionArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if string(pArgs.Path) == "" {
		return nil, fmt.Errorf("missing parameter: 'path'")
	}
	path, err := pArgs.Path.Get()
	if err != nil {
		return nil, err
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = path
	}

	err = os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("error deleting file: %w", err)
	}

	return json.Marshal(FileDeletionResult{
		Path: absolutePath,
	})
}
