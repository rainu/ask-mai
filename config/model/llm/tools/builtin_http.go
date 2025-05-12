package tools

import (
	"github.com/rainu/ask-mai/llms/tools/http"
)

type Http struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y http.CallResult    `yaml:"-"`
	Z http.CallArguments `yaml:"-"`
}

func NewHttp() Http {
	return Http{
		Approval: ApprovalNever,
	}
}

func (h Http) AsFunctionDefinition() *FunctionDefinition {
	if h.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:        "callHttp",
		Approval:    h.Approval,
		Description: http.CallDefinition.Description,
		Parameters:  http.CallDefinition.Parameter,
		CommandFn:   http.CallDefinition.Function,
	}
}
