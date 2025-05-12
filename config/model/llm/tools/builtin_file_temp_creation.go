package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type FileTempCreation struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileTempCreationResult    `yaml:"-"`
	Z file.FileTempCreationArguments `yaml:"-"`
}

func NewFileTempCreation() FileTempCreation {
	return FileTempCreation{
		Approval: ApprovalNever,
	}
}

func (f FileTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createTempFile",
		Approval:    f.Approval,
		Description: file.FileTempCreationDefinition.Description,
		Parameters:  file.FileTempCreationDefinition.Parameter,
		CommandFn:   file.FileTempCreationDefinition.Function,
	}
}
