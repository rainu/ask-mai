package file

import (
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestTool_FileTempCreation(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileTempCreationTool, FileTempCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileTempCreationTool.Name

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedResult FileTempCreationResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedResult))
	defer func() {
		os.Remove(parsedResult.Path)
	}()

	assert.Contains(t, parsedResult.Path, os.TempDir())

	info, err := os.Stat(parsedResult.Path)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm(), "File permissions should be set to 0644")
}

func TestTool_FileTempCreation_WithSuffix(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileTempCreationTool, FileTempCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileTempCreationTool.Name
	req.Params.Arguments = map[string]any{
		"suffix": ".txt",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedResult FileTempCreationResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedResult))
	defer func() {
		os.Remove(parsedResult.Path)
	}()

	assert.True(t, strings.HasPrefix(parsedResult.Path, os.TempDir()))
	assert.True(t, strings.HasSuffix(parsedResult.Path, ".txt"))

	info, err := os.Stat(parsedResult.Path)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm(), "File permissions should be set to 0644")
}

func TestTool_FileTempCreation_CustomPermissions(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(FileTempCreationTool, FileTempCreationToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = FileTempCreationTool.Name
	req.Params.Arguments = map[string]any{
		"permission": "0600",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedResult FileTempCreationResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedResult))
	defer func() {
		os.Remove(parsedResult.Path)
	}()

	assert.Contains(t, parsedResult.Path, os.TempDir())

	info, err := os.Stat(parsedResult.Path)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "File permissions should be set to 0600")
}
