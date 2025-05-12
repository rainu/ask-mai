package tools

import (
	"github.com/rainu/ask-mai/llms/tools/file"
)

type Stats struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y file.StatsResult    `yaml:"-"`
	Z file.StatsArguments `yaml:"-"`
}

func NewStats() Stats {
	return Stats{
		Approval: ApprovalNever,
	}
}

func (f Stats) AsFunctionDefinition() *FunctionDefinition {
	if f.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "getStats",
		Approval:    f.Approval,
		Description: file.StatsDefinition.Description,
		Parameters:  file.StatsDefinition.Parameter,
		CommandFn:   file.StatsDefinition.Function,
	}
}
