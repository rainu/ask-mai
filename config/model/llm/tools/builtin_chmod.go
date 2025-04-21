package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type ChangeMode struct {
	Disable  bool   `config:"disable" yaml:"disable" usage:"Disable tool"`
	Approval string `config:"approval" yaml:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.ChangeModeResult    `config:"-" yaml:"-"`
	Z file.ChangeModeArguments `config:"-" yaml:"-"`
}

func NewChangeMode() ChangeMode {
	return ChangeMode{
		Approval: ApprovalAlways,
	}
}

func (f ChangeMode) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "changeMode",
		Approval:    f.Approval,
		Description: file.ChangeModeDefinition.Description,
		Parameters:  file.ChangeModeDefinition.Parameter,
		CommandFn:   file.ChangeModeDefinition.Function,
	}
}
