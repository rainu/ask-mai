package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	it "github.com/rainu/ask-mai/config/model/llm/tools"
)

const McpPrefix = it.BuiltInPrefix + "_"

type Config struct {
	CommandServer []Command `config:"command" yaml:"command" usage:"CommandServer"`
	HttpServer    []Http    `config:"http" yaml:"http" usage:"HTTPServer"`
}

type Tool struct {
	mcp.ToolRetType
	Transport interface {
		GetTransport() transport.Transport
	}
}

func (c *Config) Validate() error {
	for _, cmd := range c.CommandServer {
		if ve := cmd.Validate(); ve != nil {
			return ve
		}
	}
	for _, http := range c.HttpServer {
		if ve := http.Validate(); ve != nil {
			return ve
		}
	}

	return nil
}

func listTools(ctx context.Context, transport transport.Transport) ([]mcp.ToolRetType, error) {
	mcpClient := mcp.NewClient(transport)
	_, err := mcpClient.Initialize(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mcp client: %w", err)
	}

	resp, err := mcpClient.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	return resp.Tools, nil
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
			}
		}
	}

	return allTools, nil
}
