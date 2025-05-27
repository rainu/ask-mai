package file

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTool_DirTempCreation(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(DirectoryTempCreationTool, DirectoryTempCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = DirectoryTempCreationTool.Name

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, os.TempDir())
}
