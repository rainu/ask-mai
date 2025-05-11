package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	it "github.com/rainu/ask-mai/config/model/llm/tools"
	"sync"
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

type listToolResult struct {
	server string
	tools  []mcp.ToolRetType
	err    error
}

func MergeTools(ctx context.Context, s map[string]Server) (map[string]Tool, error) {
	result, err := ListTools(ctx, s)
	if err != nil {
		return nil, err
	}

	allTools := map[string]Tool{}
	for serverName, tools := range result {
		for toolName, tool := range tools {
			server := s[serverName]

			// to prevent naming collisions, we add a prefix to the tool name
			allTools[fmt.Sprintf("%s%s_%s", McpPrefix, serverName, toolName)] = Tool{
				ToolRetType: tool,
				Transport:   &server,
				Timeout:     *server.Timeout.Execution,
				approval:    Approval(server.Approval),
			}
		}
	}

	return allTools, nil
}

func ListTools(ctx context.Context, s map[string]Server) (map[string]map[string]mcp.ToolRetType, error) {
	allTools := map[string]map[string]mcp.ToolRetType{}
	resultChan := make(chan listToolResult)
	wg := sync.WaitGroup{}

	// call ListTools for each server in parallel
	for name := range s {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			server := s[name]
			result := listToolResult{server: name}
			result.tools, result.err = server.ListTools(ctx)
			if result.err != nil {
				result.err = fmt.Errorf("failed to list tools for server %s: %w", name, result.err)
			}

			resultChan <- result
		}(name)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// collect results from all servers
	errors := make([]error, 0, len(s))
	for result := range resultChan {
		if result.err != nil {
			errors = append(errors, result.err)
			continue
		}
		allTools[result.server] = map[string]mcp.ToolRetType{}
		for _, tool := range result.tools {
			allTools[result.server][tool.Name] = tool
		}
	}
	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("failed to list tools for some servers: %v", errors)
	}

	return allTools, err
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
