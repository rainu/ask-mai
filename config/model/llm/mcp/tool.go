package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
)

type Tool struct {
	mcp.ToolRetType
	Transport interface {
		GetTransport() transport.Transport
	}
	approval Approval
}

func (c *Config) ListTools(ctx context.Context) (map[string]Tool, error) {
	allTools := map[string]Tool{}
	for i, cmd := range c.CommandServer {
		tools, err := cmd.ListTools(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for command server: %w", err)
		}

		for _, tool := range tools {
			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s%d_%s", McpPrefix, i, tool.Name)] = Tool{
				ToolRetType: tool,
				Transport:   &cmd,
				approval:    Approval(cmd.Approval),
			}
		}
	}

	for i, http := range c.HttpServer {
		tools, err := http.ListTools(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for http server: %w", err)
		}

		for _, tool := range tools {
			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s%d_%s", McpPrefix, i+len(c.CommandServer), tool.Name)] = Tool{
				ToolRetType: tool,
				Transport:   &http,
				approval:    Approval(http.Approval),
			}
		}
	}

	return allTools, nil
}

func listTools(ctx context.Context, transport transport.Transport) ([]mcp.ToolRetType, error) {
	mcpClient := mcp.NewClient(transport)
	_, err := mcpClient.Initialize(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mcp client: %w", err)
	}

	c := ""
	resp, err := mcpClient.ListTools(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	return resp.Tools, nil
}

func (t *Tool) NeedApproval(ctx context.Context, jsonArgs string) bool {
	return t.approval.NeedsApproval(ctx, jsonArgs, &t.ToolRetType)
}
