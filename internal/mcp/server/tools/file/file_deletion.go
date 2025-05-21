package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/mcp/server/tools"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type FileDeletionArguments struct {
	Path Path `json:"path"`
}

type FileDeletionResult struct {
	Path string `json:"path"`
}

var FileDeletionDefinition = tools.BuiltinDefinition{
	Description: "Delete a file on the user's system.",
	Parameter: map[string]any{
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
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
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
	},
}

var FileDeletionTool = mcp.NewTool("deleteFile",
	mcp.WithDescription("Delete a file on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the file to delete. Use '~' as placeholder for the user's home directory."),
	),
)

var FileDeletionToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs FileDeletionArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
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

	raw, err := json.Marshal(FileDeletionResult{
		Path: absolutePath,
	})
	return mcp.NewToolResultText(string(raw)), err
}
