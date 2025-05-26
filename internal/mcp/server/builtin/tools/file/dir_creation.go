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

type DirectoryCreationArguments struct {
	Path       Path       `json:"path"`
	Permission Permission `json:"permission"`
}

type DirectoryCreationResult struct {
	Path string `json:"path"`
}

var DirectoryCreationTool = mcp.NewTool("createDirectory",
	mcp.WithDescription("Creates a new directory (including all missing parent directories) on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the directory to create. Use '~' as placeholder for the user's home directory."),
	),
	mcp.WithString("permission",
		mcp.Description("The permission of the directory. Default is 0755."),
	),
)

var DirectoryCreationToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs DirectoryCreationArguments

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

	// Check if directory already exists
	dirInfo, dirErr := os.Stat(path)
	if dirErr == nil {
		if !dirInfo.IsDir() {
			return nil, fmt.Errorf("path exists but is a file: %s", path)
		}
		return nil, fmt.Errorf("directory already exists: %s", path)
	}

	perm, err := pArgs.Permission.Get(os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(path, perm)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn("Error getting absolute path!", "error", err)
		absolutePath = path
	}

	raw, err := json.Marshal(DirectoryCreationResult{
		Path: absolutePath,
	})
	return mcp.NewToolResultText(string(raw)), err
}
