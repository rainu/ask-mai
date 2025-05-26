package mcp

import (
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/mcp/client"
	mcpServer "github.com/rainu/ask-mai/internal/mcp/server/builtin"
	"github.com/stretchr/testify/assert"
	"testing"
)

type builtinTools struct {
	builtin.BuiltIns
}

func (b *builtinTools) GetTransport() (transport.Interface, error) {
	return transport.NewInProcessTransport(mcpServer.NewServer("", b.BuiltIns)), nil
}

func (b *builtinTools) GetTimeouts() client.Timeouts {
	return client.Timeouts{}
}

func TestClientServer(t *testing.T) {
	c, err := client.GetClient(t.Context(), &builtinTools{})
	assert.NoError(t, err)

	req := mcp.CallToolRequest{}
	req.Params.Name = "getSystemTime"

	resp, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Content, 1)
}
