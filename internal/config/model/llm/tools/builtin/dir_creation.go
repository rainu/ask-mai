package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/file"
)

type DirectoryCreation struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryCreationResult    `yaml:"-"`
	Z file.DirectoryCreationArguments `yaml:"-"`
}

func (c *DirectoryCreation) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
