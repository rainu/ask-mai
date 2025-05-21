package system

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/mcp/server/tools"
	"time"
)

var SystemTimeDefinition = tools.BuiltinDefinition{
	Description: "Get the current system time.",
	Parameter: map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		return []byte(time.Now().String()), nil
	},
}

var SystemTimeTool = mcp.NewTool("getSystemTime",
	mcp.WithDescription("Get the current system time."),
)
var SystemTimeToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(time.Now().String()), nil
}
