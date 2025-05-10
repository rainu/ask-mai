package controller

import (
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/rainu/ask-mai/config/model"
	internalMCP "github.com/rainu/ask-mai/config/model/llm/mcp"
	"github.com/rainu/ask-mai/config/model/llm/tools"
)

type ApplicationConfig struct {
	Config        model.Config
	ActiveProfile model.Profile
}

func (c *Controller) GetApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		Config:        *c.appConfig,
		ActiveProfile: *c.getProfile(),
	}
}

func (c *Controller) GetDebugConfig() model.DebugConfig {
	return c.appConfig.DebugConfig
}

func (c *Controller) GetAvailableProfiles() map[string]model.ProfileMeta {
	return c.appConfig.GetProfiles()
}

func (c *Controller) SetActiveProfile(profileName string) model.Profile {
	c.appConfig.ActiveProfile = profileName
	return *c.getProfile()
}

func (c *Controller) SetBuiltinTools(config tools.BuiltIns) {
	c.getProfile().LLM.Tools.BuiltInTools = config
}

func (c *Controller) SetMcpTools(config internalMCP.Config) {
	c.getProfile().LLM.McpServer = config
}

func (c *Controller) ListMcpCommandTools() ([][]mcp.ToolRetType, error) {
	result := [][]mcp.ToolRetType{}

	for i, cmd := range c.getProfile().LLM.McpServer.CommandServer {
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

	for i, http := range c.getProfile().LLM.McpServer.HttpServer {
		r, err := http.ListAllTools(c.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for http server #%d: %w", i, err)
		}
		result = append(result, r)
	}

	return result, nil
}

func (c *Controller) getProfile() *model.Profile {
	return c.appConfig.GetActiveProfile()
}
