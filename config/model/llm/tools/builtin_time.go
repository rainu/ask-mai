package tools

import (
	"github.com/rainu/ask-mai/llms/tools"
)

type SystemTime struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"Disable tool"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`
}

func NewSystemTime() SystemTime {
	return SystemTime{
		Approval: ApprovalNever,
	}
}

func (s SystemTime) AsFunctionDefinition() *FunctionDefinition {
	if s.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getSystemTime",
		Approval:    s.Approval,
		Description: tools.SystemTimeDefinition.Description,
		Parameters:  tools.SystemTimeDefinition.Parameter,
		CommandFn:   tools.SystemTimeDefinition.Function,
	}
}
