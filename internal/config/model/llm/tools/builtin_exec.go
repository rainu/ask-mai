package tools

import (
	"github.com/rainu/ask-mai/internal/mcp/server/tools/command"
)

type CommandExecution struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Z command.CommandExecutionArguments `yaml:"-"`
}

func (c *CommandExecution) SetDefaults() {
	if c.Approval == "" {
		c.Approval = ApprovalAlways
	}
}

func NewCommandExecution() CommandExecution {
	return CommandExecution{
		Approval: ApprovalAlways,
	}
}

func (c CommandExecution) AsFunctionDefinition() *FunctionDefinition {
	if c.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "executeCommand",
		Approval:    c.Approval,
		Description: command.CommandExecutionDefinition.Description,
		Parameters:  command.CommandExecutionDefinition.Parameter,
		CommandFn:   command.CommandExecutionDefinition.Function,
	}
}
