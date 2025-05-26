package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/file"
)

type Stats struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.StatsResult    `yaml:"-"`
	Z file.StatsArguments `yaml:"-"`
}

func (c *Stats) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
