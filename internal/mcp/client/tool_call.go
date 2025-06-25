package client

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
	if len(req.Params.Arguments) == 0 {
		req.Params.Arguments = map[string]any{
			"_": "_", // some tools require at least one argument, so we add a dummy one
		}
	}

	execCtx, cancel := tp.GetTimeouts().ExecutionContext(ctx)
	defer cancel()

	return c.CallTool(execCtx, req)
}
