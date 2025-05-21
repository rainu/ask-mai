package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientServer(t *testing.T) {
	c, err := client.GetClient(t.Context(), &tools.BuiltIns{})
	assert.NoError(t, err)

	req := mcp.CallToolRequest{}
	req.Params.Name = "getSystemTime"

	resp, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Content, 1)
}
