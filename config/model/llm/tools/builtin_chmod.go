package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeMode struct {
	Disable    bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NoApproval bool `config:"no-approval" yaml:"no-approval" usage:"Needs no user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.ChangeModeResult    `config:"-" yaml:"-"`
	Z file.ChangeModeArguments `config:"-" yaml:"-"`
}

func (f ChangeMode) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "changeMode",
		NeedsApproval: !f.NoApproval,
		Description:   file.ChangeModeDefinition.Description,
		Parameters:    file.ChangeModeDefinition.Parameter,
		CommandFn:     file.ChangeModeDefinition.Function,
	}
}
