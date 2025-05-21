package tools

import (
	"github.com/rainu/ask-mai/internal/mcp/server/tools/file"
)

type DirectoryTempCreation struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.DirectoryTempCreationResult    `yaml:"-"`
	Z file.DirectoryTempCreationArguments `yaml:"-"`
}

func (c *DirectoryTempCreation) SetDefaults() {
	if c.Approval == "" {
		c.Approval = ApprovalNever
	}
}

func NewDirectoryTempCreation() DirectoryTempCreation {
	return DirectoryTempCreation{
		Approval: ApprovalNever,
	}
}

func (f DirectoryTempCreation) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "createTempDirectory",
		Approval:    f.Approval,
		Description: file.DirectoryTempCreationDefinition.Description,
		Parameters:  file.DirectoryTempCreationDefinition.Parameter,
		CommandFn:   file.DirectoryTempCreationDefinition.Function,
	}
}
