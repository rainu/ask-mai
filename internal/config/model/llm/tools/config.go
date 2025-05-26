package tools

import (
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/mcp"
)

type Config struct {
	BuiltIns  builtin.BuiltIns                      `yaml:"builtin,omitempty" usage:"Built-in tool "`
	McpServer map[string]mcp.Server                 `yaml:"mcpServers,omitempty" usage:"MCP server "`
	Custom    map[string]command.FunctionDefinition `yaml:"custom,omitempty" usage:"Custom tool definition "`
}

func (c *Config) Validate() error {
	for name, server := range c.McpServer {
		if ve := server.Validate(); ve != nil {
			return fmt.Errorf("invlalid mcpServer config for %s: %w", name, ve)
		}
	}

	for cmd, definition := range c.Custom {
		if definition.CommandExpr != "" {
			if ve := command.Expression(definition.CommandExpr).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = command.Expression(definition.CommandExpr).CommandFn(definition)
		} else if definition.Command != "" {
			if ve := command.Command(definition.Command).Validate(); ve != nil {
				return ve
			}
			definition.CommandFn = command.Command(definition.Command).CommandFn(definition)
		} else {
			return fmt.Errorf("Command for tool '%s' is missing", cmd)
		}
	}

	return nil
}
