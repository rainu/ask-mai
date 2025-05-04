package controller

import (
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/rainu/ask-mai/config/model"
	internalMCP "github.com/rainu/ask-mai/config/model/llm/mcp"
	"github.com/rainu/ask-mai/config/model/llm/tools"
)

func (c *Controller) GetApplicationConfig() model.Config {
	return *c.getConfig()
}

func (c *Controller) GetAvailableProfiles() map[string]model.Profile {
	return c.appConfig.GetProfiles()
}

func (c *Controller) SetActiveProfile(profileName string) model.Config {
	c.appConfig.Profile.Active = profileName
	return c.GetApplicationConfig()
}

func (c *Controller) SetBuiltinTools(config tools.BuiltIns) {
	c.getConfig().LLM.Tools.BuiltInTools = config
}

func (c *Controller) SetMcpTools(config internalMCP.Config) {
	c.getConfig().LLM.McpServer = config
}

func (c *Controller) ListMcpCommandTools() ([][]mcp.ToolRetType, error) {
	result := [][]mcp.ToolRetType{}

	for i, cmd := range c.getConfig().LLM.McpServer.CommandServer {
		r, err := cmd.ListAllTools(c.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for command server #%d: %w", i, err)
		}
		result = append(result, r)
	}

	return result, nil
}

func (c *Controller) ListMcpHttpTools() ([][]mcp.ToolRetType, error) {
	result := [][]mcp.ToolRetType{}

	for i, http := range c.getConfig().LLM.McpServer.HttpServer {
		r, err := http.ListAllTools(c.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for http server #%d: %w", i, err)
		}
		result = append(result, r)
	}

	return result, nil
}

func (c *Controller) getConfig() *model.Config {
	return c.appConfig.GetActiveProfile()
}
