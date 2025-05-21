package mcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"sync"
)

func ListTools(ctx context.Context, tp Transporter) ([]mcp.Tool, error) {
	c, err := GetClient(ctx, tp)
	if err != nil {
		return nil, err
	}

	listCtx, cancel := tp.GetTimeouts().ListContext(ctx)
	defer cancel()

	resp, err := c.ListTools(listCtx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	return resp.Tools, nil
}

type listToolResult struct {
	server string
	tools  []mcp.Tool
	err    error
}

func ListAllTools(ctx context.Context, s map[string]Transporter) (map[string]map[string]mcp.Tool, error) {
	allTools := map[string]map[string]mcp.Tool{}
	resultChan := make(chan listToolResult)
	wg := sync.WaitGroup{}

	// call ListTools for each server in parallel
	for name := range s {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			server := s[name]
			result := listToolResult{server: name}
			result.tools, result.err = ListTools(ctx, server)
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
		allTools[result.server] = map[string]mcp.Tool{}
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
