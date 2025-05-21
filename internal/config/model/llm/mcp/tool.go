package mcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	it "github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"slices"
)

const McpPrefix = it.BuiltInPrefix + "_"

type Tool struct {
	mcp.Tool
	Transporter client.Transporter
	approval    Approval
}

func MergeTools(ctx context.Context, s map[string]Server) (map[string]Tool, error) {
	tp := make(map[string]client.Transporter, len(s))
	for name, server := range s {
		tp[name] = &server
	}

	result, err := client.ListAllTools(ctx, tp)
	if err != nil {
		return nil, err
	}

	allTools := map[string]Tool{}
	idx := 0
	for serverName, tools := range result {
		for toolName, tool := range tools {
			server := s[serverName]

			if slices.Contains(server.Exclude, tool.Name) {
				// skip excluded tools
				continue
			}

			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s%d%s", McpPrefix, idx, toolName)] = Tool{
				Tool:        tool,
				Transporter: &server,
				approval:    Approval(server.Approval),
			}
		}

		idx++
	}

	return allTools, nil
}

func (t *Tool) NeedApproval(ctx context.Context, jsonArgs string) bool {
	return t.approval.NeedsApproval(ctx, jsonArgs, &t.Tool)
}
