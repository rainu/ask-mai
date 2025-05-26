package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type DirectoryDeletionArguments struct {
	Path Path `json:"path"`
}

type DirectoryDeletionResult struct {
	Path string `json:"path"`
}

var DirectoryDeletionTool = mcp.NewTool("deleteDirectory",
	mcp.WithDescription("Delete a directory (including all files ans subdirectories) on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the directory to delete. Use '~' as placeholder for the user's home directory."),
	),
)

var DirectoryDeletionToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs DirectoryDeletionArguments

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

	err = os.RemoveAll(path)
	if err != nil {
		return nil, fmt.Errorf("error deleting directory: %w", err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = path
	}

	raw, err := json.Marshal(DirectoryDeletionResult{
		Path: absolutePath,
	})
	return mcp.NewToolResultText(string(raw)), err
}
