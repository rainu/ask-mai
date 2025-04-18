package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type SystemTime struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`
}

func (s SystemTime) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "getSystemTime",
		NeedsApproval: s.NeedsApproval,
		Description:   tools.SystemTimeDefinition.Description,
		Parameters:    tools.SystemTimeDefinition.Parameter,
		CommandFn:     tools.SystemTimeDefinition.Function,
	}
}
