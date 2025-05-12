package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type Environment struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y tools.EnvironmentResult    `yaml:"-"`
	Z tools.EnvironmentArguments `yaml:"-"`
}

func NewEnvironment() Environment {
	return Environment{
		Approval: ApprovalNever,
	}
}

func (f Environment) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getEnvironment",
		Approval:    f.Approval,
		Description: tools.EnvironmentDefinition.Description,
		Parameters:  tools.EnvironmentDefinition.Parameter,
		CommandFn:   tools.EnvironmentDefinition.Function,
	}
}
