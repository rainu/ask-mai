package builtin

import (
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/http"
)

type Http struct {
	Disable  bool   `yaml:"disable,omitempty" usage:"disable"`
	Approval string `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute this tool"`

	//only for wails to generate TypeScript types
	Y http.CallResult    `yaml:"-"`
	Z http.CallArguments `yaml:"-"`
}

func (c *Http) SetDefaults() {
	if c.Approval == "" {
		c.Approval = approval.Never
	}
}
