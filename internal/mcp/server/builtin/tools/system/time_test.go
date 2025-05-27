package system

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTool_Time(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(SystemTimeTool, SystemTimeToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = SystemTimeTool.Name

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	text := res.Content[0].(mcp.TextContent).Text
	now := time.Now().String()[0:10]
	assert.Contains(t, text, now)
}
