package tools

import (
	"github.com/rainu/ask-mai/llms/tools/command"
)

type CommandExecution struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Z command.CommandExecutionArguments `config:"-" yaml:"-"`
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
