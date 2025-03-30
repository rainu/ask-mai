package tools

import (
	"github.com/rainu/ask-mai/llms/tools/http"
)

type Http struct {
	Disable       bool `config:"disable" yaml:"disable" usage:"Disable tool"`
	NeedsApproval bool `config:"approval" yaml:"approval" usage:"Needs user approval to be executed"`

	//only for wails to generate TypeScript types
	Y http.CallResult    `config:"-" yaml:"-"`
	Z http.CallArguments `config:"-" yaml:"-"`
}

func (h Http) AsFunctionDefinition() *FunctionDefinition {
	if h.Disable {
		return nil
	}

	return &FunctionDefinition{
		Name:          "callHttp",
		NeedsApproval: h.NeedsApproval,
		Description:   http.CallDefinition.Description,
		Parameters:    http.CallDefinition.Parameter,
		CommandFn:     http.CallDefinition.Function,
	}
}
