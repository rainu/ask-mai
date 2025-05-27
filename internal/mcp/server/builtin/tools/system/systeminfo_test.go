package system

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestTool_Info(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(SystemInfoTool, SystemInfoToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = SystemInfoTool.Name

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	text := res.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, runtime.GOOS)
	assert.Contains(t, text, runtime.GOARCH)
}
