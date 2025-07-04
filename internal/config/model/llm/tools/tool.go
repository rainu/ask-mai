package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"slices"
)

const (
	ServerNameBuiltin = "_b"
	ServerNameCustom  = "_c"
)

type Tool struct {
	mcp.Tool
	Transporter client.Transporter
	approval    approval.Approval
}

func (c *Config) GetBuiltinTransporter() client.Transporter {
	return &builtinTools{BuiltIns: c.BuiltIns}
}

func (c *Config) GetCustomTransporter() client.Transporter {
	return &customTools{config: c.Custom}
}

func (c *Config) GetTools(ctx context.Context) (map[string]Tool, error) {
	tp := map[string]client.Transporter{}
	for name, server := range c.McpServer {
		tp[name] = &server
	}
	tp[ServerNameBuiltin] = c.GetBuiltinTransporter()
	tp[ServerNameCustom] = c.GetCustomTransporter()

	result, err := client.ListAllTools(ctx, tp)
	if err != nil {
		return nil, err
	}

	allTools := map[string]Tool{}

	for serverName, tools := range result {
		for toolName, tool := range tools {
			var toolApproval approval.Approval
			if serverName == ServerNameBuiltin {
				toolApproval = approval.Approval(c.BuiltIns.GetApprovalFor(toolName))
			} else if serverName == ServerNameCustom {
				toolApproval = approval.Approval(c.Custom[toolName].Approval)
			} else {
				server := c.McpServer[serverName]
				if slices.Contains(server.Exclude, tool.Name) {
					// skip excluded tools
					continue
				}

				toolApproval = approval.Approval(server.Approval)
			}

			allTools[GetUniqToolName(serverName, toolName)] = Tool{
				Tool:        tool,
				Transporter: tp[serverName],
				approval:    toolApproval,
			}
		}
	}

	return allTools, nil
}

func GetUniqToolName(serverName, toolName string) string {
	// to prevent naming collisions, we add a prefix to the tool name
	return fmt.Sprintf("%s_%s", serverName, toolName)
}

func (t *Tool) NeedApproval(ctx context.Context, jsonArgs string) bool {
	return t.approval.NeedsApproval(ctx, jsonArgs, &t.Tool)
}
