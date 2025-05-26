package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/system"
)

type Environment struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y system.EnvironmentResult    `yaml:"-"`
	Z system.EnvironmentArguments `yaml:"-"`
}

func (c *Environment) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
