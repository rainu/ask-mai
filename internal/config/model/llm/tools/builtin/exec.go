package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/command"
)

type CommandExecution struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Z command.CommandExecutionArguments `yaml:"-"`
}

func (c *CommandExecution) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Always
	}
}
