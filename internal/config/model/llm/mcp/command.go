package mcp

import (
	"fmt"
	"github.com/mark3labs/mcp-go/client/transport"
)

type Command struct {
	Name        string            `yaml:"command,omitempty" usage:"[Command] Name of the command"`
	Arguments   []string          `yaml:"args,omitempty" usage:"[Command] Arguments to pass"`
	Environment map[string]string `yaml:"env,omitempty" usage:"[Command] Environment variables to pass"`

	env []string `yaml:"-"`
}

func (c *Command) Validate() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("command name is required")
	}
	return nil
}

func (c *Command) Env() []string {
	if c.env == nil {
		c.env = []string{}
		for key, value := range c.Environment {
			c.env = append(c.env, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return c.env
}

func (c *Command) GetTransport() (transport.Interface, error) {
	return transport.NewStdio(c.Name, c.Env(), c.Arguments...), nil
}
