package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeTimes struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.ChangeTimesResult    `config:"-" yaml:"-"`
	Z file.ChangeTimesArguments `config:"-" yaml:"-"`
}

func NewChangeTimes() ChangeTimes {
	return ChangeTimes{
		Approval: ApprovalNever,
	}
}

func (f ChangeTimes) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "changeTimes",
		Approval:    f.Approval,
		Description: file.ChangeTimesDefinition.Description,
		Parameters:  file.ChangeTimesDefinition.Parameter,
		CommandFn:   file.ChangeTimesDefinition.Function,
	}
}
