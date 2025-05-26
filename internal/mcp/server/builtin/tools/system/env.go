package system

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"strings"
)

type EnvironmentArguments struct {
}

type EnvironmentResult struct {
	Environment map[string]string `json:"env"`
}

var EnvironmentTool = mcp.NewTool("getEnvironment",
	mcp.WithDescription("Get all environment variables of the user's system."),
)

var EnvironmentToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := EnvironmentResult{
		Environment: map[string]string{},
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		result.Environment[pair[0]] = pair[1]
	}

	raw, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(raw)), nil
}
