package tools

import (
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	"github.com/rainu/ask-mai/internal/mcp/client"
	mcpServer "github.com/rainu/ask-mai/internal/mcp/server/custom"
)

type customTools struct {
	config map[string]command.FunctionDefinition
}

func (c *customTools) GetTransport() (transport.Interface, error) {
	return transport.NewInProcessTransport(mcpServer.NewServer("", c.config)), nil
}

func (c *customTools) GetTimeouts() client.Timeouts {
	return client.Timeouts{} // No timeouts for custom tools
}
