package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeOwner struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.ChangeOwnerResult    `config:"-" yaml:"-"`
	Z file.ChangeOwnerArguments `config:"-" yaml:"-"`
}

func NewChangeOwner() ChangeOwner {
	return ChangeOwner{
		Approval: ApprovalAlways,
	}
}

func (f ChangeOwner) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "changeOwner",
		Approval:    f.Approval,
		Description: file.ChangeOwnerDefinition.Description,
		Parameters:  file.ChangeOwnerDefinition.Parameter,
		CommandFn:   file.ChangeOwnerDefinition.Function,
	}
}
