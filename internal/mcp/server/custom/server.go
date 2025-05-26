package custom

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
)

func NewServer(version string, cfg map[string]command.FunctionDefinition) *server.MCPServer {
	s := server.NewMCPServer(
		"ask-mai",
		version,
		server.WithToolCapabilities(false),
	)

	for name, definition := range cfg {
		t := mcp.Tool{
			Name:        name,
			Description: definition.Description,
			InputSchema: definition.Parameters,
		}
		s.AddTool(t, handlerFor(definition))
	}

	return s
}

func handlerFor(definition command.FunctionDefinition) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw, err := json.Marshal(request.Params.Arguments)
		if err != nil {
			return nil, err
		}

		rawResult, err := definition.CommandFn(ctx, string(raw))
		return mcp.NewToolResultText(string(rawResult)), err
	}
}
