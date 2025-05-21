package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func CallTool(ctx context.Context, tp Transporter, toolName string, argsAsJson string) (*mcp.CallToolResult, error) {
	c, err := GetClient(ctx, tp)
	if err != nil {
		return nil, err
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName

	err = json.Unmarshal([]byte(argsAsJson), &req.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool call arguments: %w", err)
	}

	execCtx, cancel := tp.GetTimeouts().ExecutionContext(ctx)
	defer cancel()

	return c.CallTool(execCtx, req)
}
