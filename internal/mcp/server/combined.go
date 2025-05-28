package server

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	bServer "github.com/rainu/ask-mai/internal/mcp/server/builtin"
	cServer "github.com/rainu/ask-mai/internal/mcp/server/custom"
)

func NewServer(version string, bConfig builtin.BuiltIns, cConfig map[string]command.FunctionDefinition) *server.MCPServer {
	s := server.NewMCPServer(
		"ask-mai",
		version,
		server.WithToolCapabilities(false),
	)
	bServer.AddTools(s, bConfig)
	cServer.AddTools(s, cConfig)

	return s
}
