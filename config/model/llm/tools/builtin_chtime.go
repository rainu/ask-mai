package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeTimes struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"Disable tool"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.ChangeTimesResult    `yaml:"-"`
	Z file.ChangeTimesArguments `yaml:"-"`
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
