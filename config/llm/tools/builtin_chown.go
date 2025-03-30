package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeOwner struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.ChangeOwnerResult    `config:"-" yaml:"-"`
	Z file.ChangeOwnerArguments `config:"-" yaml:"-"`
}

func (f ChangeOwner) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "changeOwner",
		NeedsApproval: !f.NoApproval,
		Description:   file.ChangeOwnerDefinition.Description,
		Parameters:    file.ChangeOwnerDefinition.Parameter,
		CommandFn:     file.ChangeOwnerDefinition.Function,
	}
}
