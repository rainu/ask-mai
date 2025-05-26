package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"os"
)

type ChangeModeArguments struct {
	Path       Path       `json:"path"`
	Permission Permission `json:"permission"`
}

type ChangeModeResult struct {
}

var ChangeModeTool = mcp.NewTool("changeMode",
	mcp.WithDescription("Changes the mode of file or directory on the user's system."),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The path to the file or directory to change the mode for. Use '~' as placeholder for the user's home directory."),
	),
	mcp.WithString("permission",
		mcp.Required(),
		mcp.Description("The permission of the file or directory."),
	),
)

var ChangeModeToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs ChangeModeArguments

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

	if string(pArgs.Permission) == "" {
		return nil, fmt.Errorf("missing parameter: 'permission'")
	}
	perm, err := pArgs.Permission.Get(os.FileMode(0000))
	if err != nil {
		return nil, err
	}

	err = os.Chmod(path, perm)
	if err != nil {
		return nil, fmt.Errorf("error changing mode: %w", err)
	}

	return mcp.NewToolResultText(""), nil
}
