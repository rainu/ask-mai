package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeTimes struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y file.ChangeTimesResult    `config:"-" yaml:"-"`
	Z file.ChangeTimesArguments `config:"-" yaml:"-"`
}

func (f ChangeTimes) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "changeTimes",
		NeedsApproval: f.NeedsApproval,
		Description:   file.ChangeTimesDefinition.Description,
		Parameters:    file.ChangeTimesDefinition.Parameter,
		CommandFn:     file.ChangeTimesDefinition.Function,
	}
}
