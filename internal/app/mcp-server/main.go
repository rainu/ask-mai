package mcp_server

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	imcpServer "github.com/rainu/ask-mai/internal/mcp/server"
	"os"
)

type Args struct {
	VersionLine string
}

func Main(args Args) int {
	s := imcpServer.NewServer(args.VersionLine, imcpServer.Options{})
	err := server.ServeStdio(s)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
