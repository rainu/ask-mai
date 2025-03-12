package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type FileCreation struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `yaml:"approval" json:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y FileCreationResult    `config:"-"`
	Z FileCreationArguments `config:"-"`
}

func (f FileCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createFile",
		Description: "Creates a file on the user's system.",
		CommandFn:   f.Command,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path to the file to create.",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "The content of the file.",
				},
				"append": map[string]any{
					"type":        "boolean",
					"description": "Append the content to the file. Default is false.",
				},
				"permission": map[string]any{
					"type":        "string",
					"description": "The permission of the file. Default is 0644.",
				},
			},
			"additionalProperties": false,
			"required":             []string{"path"},
		},
		NeedsApproval: f.NeedsApproval,
	}
}

type FileCreationArguments struct {
	Path       string `json:"path"`
	Content    string `json:"content"`
	Append     bool   `json:"append"`
	Permission string `json:"permission"`
}

type FileCreationResult struct {
	Path string `json:"path"`
	Size int    `json:"size"`
}

func (f FileCreation) Command(ctx context.Context, jsonArguments string) ([]byte, error) {
	var pArgs FileCreationArguments
	err := json.Unmarshal([]byte(jsonArguments), &pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if pArgs.Path == "" {
		return nil, fmt.Errorf("missing parameter: 'path'")
	}

	flag := os.O_WRONLY | os.O_CREATE
	if pArgs.Append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	perm := os.FileMode(0644)
	if pArgs.Permission != "" {
		pi, pe := strconv.ParseInt(pArgs.Permission, 8, 32)
		if pe != nil {
			return nil, fmt.Errorf("error parsing permissions: %w", pe)
		}
		perm = os.FileMode(pi)
	}

	file, err := os.OpenFile(pArgs.Path, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	absolutePath, err := filepath.Abs(file.Name())
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = file.Name()
	}

	s, err := file.WriteString(pArgs.Content)
	return json.Marshal(FileCreationResult{
		Path: absolutePath,
		Size: s,
	})
}
