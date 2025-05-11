package controller

import (
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/rainu/ask-mai/config/model"
	internalMcp "github.com/rainu/ask-mai/config/model/llm/mcp"
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
	c.getProfile().LLM.Tools.BuiltInTools = &config
}

func (c *Controller) SetMcpTools(config map[string]internalMcp.Server) {
	c.getProfile().LLM.McpServer = config
}

func (c *Controller) ListMcpTools() (map[string][]mcp.ToolRetType, error) {
	result := map[string][]mcp.ToolRetType{}

	for name, server := range c.getProfile().LLM.McpServer {
		var err error
		result[name], err = server.ListAllTools(c.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for command server %s: %w", name, err)
		}
	}

	return result, nil
}

func (c *Controller) getProfile() *model.Profile {
	return c.appConfig.GetActiveProfile()
}
