package main

import (
	"embed"
	mcp_server "github.com/rainu/ask-mai/internal/app/mcp-server"
	"github.com/rainu/ask-mai/internal/app/ui"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"os"
	"path"
	"strings"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	defer client.Close()

	var rc int

	mode := getMode()

	if mode == "ask-mai-mcp-server" {
		rc = mcp_server.Main(mcp_server.Args{
			VersionLine: versionLine(),
		})
	} else {
		rc = ui.Main(ui.Args{
			Assets:      assets,
			Icon:        icon,
			VersionLine: versionLine(),
		})
	}

	os.Exit(rc)
}

func getMode() string {
	mode := path.Base(os.Args[0])
	if len(os.Args) >= 2 {
		if !strings.HasPrefix(os.Args[1], "-") {
			mode = os.Args[1]

			// remove the mode from the args
			os.Args = append(os.Args[:1], os.Args[2:]...)
		}
	}
	return mode
}
