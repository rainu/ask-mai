package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/file"
)

type FileAppending struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileAppendingResult    `yaml:"-"`
	Z file.FileAppendingArguments `yaml:"-"`
}

func (c *FileAppending) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Always
	}
}
