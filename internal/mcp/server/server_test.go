package server

import (
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	testServer := NewServer("test", Options{})

	c := client.NewClient(transport.NewInProcessTransport(testServer))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	assert.NoError(t, err)

	result, err := c.ListTools(t.Context(), mcp.ListToolsRequest{})
	assert.NoError(t, err)
	assert.Len(t, result.Tools, 17)
}
