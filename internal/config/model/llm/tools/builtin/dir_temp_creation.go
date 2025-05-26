package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/file"
)

type DirectoryTempCreation struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryTempCreationResult    `yaml:"-"`
	Z file.DirectoryTempCreationArguments `yaml:"-"`
}

func (c *DirectoryTempCreation) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
