package tools

import (
	"github.com/rainu/ask-mai/internal/llms/tools/file"
)

type FileDeletion struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.FileDeletionResult    `yaml:"-"`
	Z file.FileDeletionArguments `yaml:"-"`
}

func NewFileDeletion() FileDeletion {
	return FileDeletion{
		Approval: ApprovalAlways,
	}
}

func (f FileDeletion) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "deleteFile",
		Approval:    f.Approval,
		Description: file.FileDeletionDefinition.Description,
		Parameters:  file.FileDeletionDefinition.Parameter,
		CommandFn:   file.FileDeletionDefinition.Function,
	}
}
