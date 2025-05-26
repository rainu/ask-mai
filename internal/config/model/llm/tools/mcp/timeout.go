package mcp

import (
	"github.com/rainu/go-yacl"
	"time"
)

type Timeout struct {
	Init      *time.Duration `yaml:"init,omitempty" usage:"Timeout for initializing. Example: \"5s\" (5 seconds)"`
	List      *time.Duration `yaml:"list,omitempty" usage:"Timeout for listing tools. Example: \"5s\" (5 seconds)"`
	Execution *time.Duration `yaml:"execution,omitempty" usage:"Timeout for executing tools. Example: \"1m\" (1 minute)"`
}

func (t *Timeout) SetDefaults() {
	if t.Init == nil {
		t.Init = yacl.P(5 * time.Second)
	}
	if t.List == nil {
		t.List = yacl.P(5 * time.Second)
	}
	if t.Execution == nil {
		t.Execution = yacl.P(1 * time.Minute)
	}
}
