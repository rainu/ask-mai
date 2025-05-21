package controller

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model"
	configMcp "github.com/rainu/ask-mai/internal/config/model/llm/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	internalMcp "github.com/rainu/ask-mai/internal/llms/tools/mcp"
)

type ApplicationConfig struct {
	Config        model.Config
	ActiveProfile model.Profile

	//only for wails to generate TypeScript types
	Z mcp.Tool
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

func (c *Controller) SetMcpTools(config map[string]configMcp.Server) {
	c.getProfile().LLM.McpServer = config
}

func (c *Controller) ListMcpTools() (map[string][]mcp.Tool, error) {
	result := map[string][]mcp.Tool{}

	for name, server := range c.getProfile().LLM.McpServer {
		var err error
		result[name], err = internalMcp.ListTools(c.ctx, &server)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for command server %s: %w", name, err)
		}
	}

	return result, nil
}

func (c *Controller) getProfile() *model.Profile {
	return c.appConfig.GetActiveProfile()
}
