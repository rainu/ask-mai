package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	commitHash = "9295DFD720"
	branch     = ""
	tag        = "dev"
	built      = "629579700000"
)

func versionLine() string {
	version := tag
	if branch != "" {
		version = branch
	}
	builtTime := time.UnixMilli(0)
	if iBuilt, _ := strconv.ParseInt(built, 10, 64); iBuilt > 0 {
		builtTime = time.UnixMilli(iBuilt)
	}

	return fmt.Sprintf("%s (#%s - %s) - https://github.com/rainu/ask-mai/tree/%s",
		version,
		commitHash[:6],
		builtTime.UTC().Format(time.RFC3339),
		commitHash[:6],
	)
}
