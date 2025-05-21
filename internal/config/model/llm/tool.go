package llm

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	iMcp "github.com/rainu/ask-mai/internal/config/model/llm/mcp"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"slices"
)

const (
	builtinServerName = ""
)

type Tool struct {
	mcp.Tool
	Transporter client.Transporter
	approval    iMcp.Approval
}

func (c *LLMConfig) GetTools(ctx context.Context) (map[string]Tool, error) {
	tp := make(map[string]client.Transporter, len(c.McpServer)+1)
	for name, server := range c.McpServer {
		tp[name] = &server
	}
	tp[builtinServerName] = &c.McpBuiltin

	result, err := client.ListAllTools(ctx, tp)
	if err != nil {
		return nil, err
	}

	allTools := map[string]Tool{}

	for serverName, tools := range result {
		for toolName, tool := range tools {
			var approval iMcp.Approval
			if serverName == builtinServerName {
				approval = iMcp.Approval(c.McpBuiltin.GetApprovalFor(toolName))
			} else {
				server := c.McpServer[serverName]
				if slices.Contains(server.Exclude, tool.Name) {
					// skip excluded tools
					continue
				}

				approval = iMcp.Approval(server.Approval)
			}

			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s_%s", serverName, toolName)] = Tool{
				Tool:        tool,
				Transporter: tp[serverName],
				approval:    approval,
			}
		}
	}

	return allTools, nil
}

func (t *Tool) NeedApproval(ctx context.Context, jsonArgs string) bool {
	return t.approval.NeedsApproval(ctx, jsonArgs, &t.Tool)
}
