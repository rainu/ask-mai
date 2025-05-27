package file

import (
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestTool_Stats(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(StatsTool, StatsToolHandler)
	})

	testFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	_, err = testFile.WriteString("First line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = StatsTool.Name
	req.Params.Arguments = map[string]any{
		"path": testFile.Name(),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedRsult StatsResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedRsult))

	assert.NotEmpty(t, parsedRsult.ModTime)
	parsedRsult.ModTime = time.Time{}

	assert.Equal(t, StatsResult{
		Path:        testFile.Name(),
		IsDirectory: false,
		IsRegular:   true,
		Permissions: "-rw-------",
		Size:        12,
		ModTime:     time.Time{},
	}, parsedRsult)
}

func TestTool_Stats_HomeResolving(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(StatsTool, StatsToolHandler)
	})

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	testFile, err := os.Create(path.Join(home, fmt.Sprintf(".ask-mai-%d", time.Now().Unix())))
	require.NoError(t, err)
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	_, err = testFile.WriteString("First line.\n")
	require.NoError(t, err)
	testFile.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = StatsTool.Name
	req.Params.Arguments = map[string]any{
		"path": path.Join("~", path.Base(testFile.Name())),
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedRsult StatsResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedRsult))

	assert.NotEmpty(t, parsedRsult.ModTime)
	parsedRsult.ModTime = time.Time{}

	assert.Equal(t, StatsResult{
		Path:        testFile.Name(),
		IsDirectory: false,
		IsRegular:   true,
		Permissions: "-rw-r--r--",
		Size:        12,
		ModTime:     time.Time{},
	}, parsedRsult)
}

func TestTool_Stats_Requirements_Path(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(StatsTool, StatsToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = StatsTool.Name
	req.Params.Arguments = map[string]any{}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "missing parameter: 'path'")
}

func TestTool_Stats_Requirements_FileDoesNotExists(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(StatsTool, StatsToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = StatsTool.Name
	req.Params.Arguments = map[string]any{
		"path": "some/path/to/file.txt",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "no such file")
}

func TestTool_Stats_Requirements_PathIsDirectory(t *testing.T) {
	c := getTestClient(t, func(s *server.MCPServer) {
		s.AddTool(StatsTool, StatsToolHandler)
	})

	req := mcp.CallToolRequest{}
	req.Params.Name = StatsTool.Name
	req.Params.Arguments = map[string]any{
		"path":    t.TempDir(),
		"content": "This is a test content.",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var parsedRsult StatsResult
	require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &parsedRsult))

	assert.NotEmpty(t, parsedRsult.ModTime)
	parsedRsult.ModTime = time.Time{}

	assert.Equal(t, StatsResult{
		Path:        req.Params.Arguments["path"].(string),
		IsDirectory: true,
		IsRegular:   false,
		Permissions: "-rwxr-xr-x",
		Size:        40,
		ModTime:     time.Time{},
	}, parsedRsult)
}
