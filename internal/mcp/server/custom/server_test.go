package custom

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer(t *testing.T) {
	testServer := NewServer("test", map[string]command.FunctionDefinition{
		"mkdir": {
			Description: "Creates a directory.",
			Parameters: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path of the directory to create.",
					},
				},
				Required: []string{"path"},
			},
			CommandFn: func(ctx context.Context, jsonArguments string) ([]byte, error) {
				return []byte("OK"), nil
			},
		},
	})

	c := client.NewClient(transport.NewInProcessTransport(testServer))

	_, err := c.Initialize(t.Context(), mcp.InitializeRequest{})
	assert.NoError(t, err)

	result, err := c.ListTools(t.Context(), mcp.ListToolsRequest{})
	assert.NoError(t, err)
	assert.Len(t, result.Tools, 1)

	assert.Equal(t, mcp.Tool{
		Name:        "mkdir",
		Description: "Creates a directory.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "The path of the directory to create.",
				},
			},
			Required: []string{"path"},
		},
	}, result.Tools[0])

	req := mcp.CallToolRequest{}
	req.Params.Name = "mkdir"
	req.Params.Arguments = map[string]any{
		"path": "/tmp/test",
	}

	res, err := c.CallTool(t.Context(), req)
	assert.NoError(t, err)

	expectedResult := mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("OK"),
		},
	}
	assert.NotNil(t, res)
	assert.Equal(t, expectedResult, *res)
}
