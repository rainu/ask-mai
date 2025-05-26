package mcp_server

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	mcpServer "github.com/rainu/ask-mai/internal/mcp/server/builtin"
	"os"
)

type Args struct {
	VersionLine string
}

func Main(args Args) int {
	s := mcpServer.NewServer(args.VersionLine, builtin.BuiltIns{})
	err := server.ServeStdio(s)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
