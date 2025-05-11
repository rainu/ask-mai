package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	it "github.com/rainu/ask-mai/config/model/llm/tools"
	"time"
)

const McpPrefix = it.BuiltInPrefix + "_"

type Tool struct {
	mcp.ToolRetType
	Transport interface {
		GetTransport() transport.Transport
	}
	approval Approval
	Timeout  time.Duration
}

func ListTools(ctx context.Context, s map[string]Server) (map[string]Tool, error) {
	allTools := map[string]Tool{}

	i := 0
	for name, server := range s {
		tools, err := server.ListTools(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for server %s: %w", name, err)
		}

		for _, tool := range tools {
			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s%d_%s", McpPrefix, i, tool.Name)] = Tool{
				ToolRetType: tool,
				Transport:   &server,
				Timeout:     *server.Timeout.Execution,
				approval:    Approval(server.Approval),
			}
		}
		i++
	}

	return allTools, nil
}

func listAllTools(ctx context.Context, transport transport.Transport) ([]mcp.ToolRetType, error) {
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

func listTools(ctx context.Context, transport transport.Transport, exclusion []string) ([]mcp.ToolRetType, error) {
	result, err := listAllTools(ctx, transport)
	if err != nil {
		return nil, err
	}

	// filter out excluded tools
	for _, exclude := range exclusion {
		for i, tool := range result {
			if tool.Name == exclude {
				result = append(result[:i], result[i+1:]...)
				break
			}
		}
	}
	return result, nil
}

func (t *Tool) NeedApproval(ctx context.Context, jsonArgs string) bool {
	return t.approval.NeedsApproval(ctx, jsonArgs, &t.ToolRetType)
}
