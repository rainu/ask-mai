package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type DirectoryDeletion struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y DirectoryDeletionResult    `config:"-" yaml:"-"`
	Z DirectoryDeletionArguments `config:"-" yaml:"-"`
}

func (f DirectoryDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "deleteDirectory",
		Description: "Delete a directory (including all files ans subdirectories) on the user's system.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path to the directory to delete. Use '~' as placeholder for the user's home directory.",
				},
			},
			"additionalProperties": false,
			"required":             []string{"path"},
		},
		NeedsApproval: !f.NoApproval,
	}
}

type DirectoryDeletionArguments struct {
	Path Path `json:"path"`
}

type DirectoryDeletionResult struct {
	Path string `json:"path"`
}

func (f DirectoryDeletion) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs DirectoryDeletionArguments
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

	err = os.RemoveAll(path)
	if err != nil {
		return nil, fmt.Errorf("error deleting directory: %w", err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = path
	}

	return json.Marshal(DirectoryDeletionResult{
		Path: absolutePath,
	})
}
