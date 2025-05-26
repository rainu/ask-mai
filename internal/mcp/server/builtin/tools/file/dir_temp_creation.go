package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"os"
)

type DirectoryTempCreationArguments struct {
}

type DirectoryTempCreationResult struct {
	Path string `json:"path"`
}

var DirectoryTempCreationTool = mcp.NewTool("createTempDirectory",
	mcp.WithDescription("Creates a new temporary directory on the user's system."),
)

var DirectoryTempCreationToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs DirectoryTempCreationArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	path, err := os.MkdirTemp("", "ask-mai.*")
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	raw, err := json.Marshal(DirectoryTempCreationResult{
		Path: path,
	})
	return mcp.NewToolResultText(string(raw)), err
}
