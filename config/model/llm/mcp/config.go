package mcp

import (
	it "github.com/rainu/ask-mai/config/model/llm/tools"
)

const McpPrefix = it.BuiltInPrefix + "_"

type Config struct {
	CommandServer []Command `config:"command" yaml:"command" usage:"CommandServer"`
	HttpServer    []Http    `config:"http" yaml:"http" usage:"HTTPServer"`
}

func (c *Config) Validate() error {
	for _, cmd := range c.CommandServer {
		if ve := cmd.Validate(); ve != nil {
			return ve
		}
	}
	for _, http := range c.HttpServer {
		if ve := http.Validate(); ve != nil {
			return ve
		}
	}

	return nil
}
