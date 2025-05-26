package builtin

import (
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	testServer := NewServer("test", builtin.BuiltIns{})

	c := client.NewClient(transport.NewInProcessTransport(testServer))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	assert.NoError(t, err)

	result, err := c.ListTools(t.Context(), mcp.ListToolsRequest{})
	assert.NoError(t, err)
	assert.Len(t, result.Tools, 17)
}
