package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/tools"
	"os"
)

type DirectoryTempCreationArguments struct {
}

type DirectoryTempCreationResult struct {
	Path string `json:"path"`
}

var DirectoryTempCreationDefinition = tools.BuiltinDefinition{
	Description: "Creates a new temporary directory on the user's system.",
	Parameter: map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
		"required":             []string{},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
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
	},
}
