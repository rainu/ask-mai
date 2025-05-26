package tools

import (
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/mcp/client"
	mcpServer "github.com/rainu/ask-mai/internal/mcp/server/builtin"
)

type builtinTools struct {
	builtin.BuiltIns
}

func (b *builtinTools) GetTransport() (transport.Interface, error) {
	return transport.NewInProcessTransport(mcpServer.NewServer("", b.BuiltIns)), nil
}

func (b *builtinTools) GetTimeouts() client.Timeouts {
	return client.Timeouts{} // No timeouts for built-in tools
}
