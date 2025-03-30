package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type Environment struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y tools.EnvironmentResult    `config:"-" yaml:"-"`
	Z tools.EnvironmentArguments `config:"-" yaml:"-"`
}

func (f Environment) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "getEnvironment",
		NeedsApproval: f.NeedsApproval,
		Description:   tools.EnvironmentDefinition.Description,
		Parameters:    tools.EnvironmentDefinition.Parameter,
		CommandFn:     tools.EnvironmentDefinition.Function,
	}
}
