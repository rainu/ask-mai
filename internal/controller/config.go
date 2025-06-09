package controller

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/builtin"
	iMcp "github.com/rainu/ask-mai/internal/config/model/llm/tools/mcp"
	"github.com/rainu/ask-mai/internal/mcp/client"
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

func (c *Controller) SetActiveProfile(profileName string) (model.Profile, error) {
	c.appConfig.ActiveProfile = profileName

	newModel, err := c.appConfig.GetActiveProfile().LLM.BuildLLM()
	if err != nil {
		return model.Profile{}, err
	}

	// close the old, and set the new
	c.aiModel.Close()
	c.aiModel = newModel

	return *c.getProfile(), nil
}

func (c *Controller) SetBuiltinTools(config builtin.BuiltIns) {
	c.getProfile().LLM.Tool.BuiltIns = config
}

func (c *Controller) SetMcpTools(config map[string]iMcp.Server) {
	c.getProfile().LLM.Tool.McpServer = config
}

func (c *Controller) ListMcpTools() (map[string][]mcp.Tool, error) {
	result := map[string][]mcp.Tool{}

	for name, server := range c.getProfile().LLM.Tool.McpServer {
		var err error
		result[name], err = client.ListTools(c.ctx, &server)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools for command server %s: %w", name, err)
		}
	}

	return result, nil
}

func (c *Controller) getProfile() *model.Profile {
	return c.appConfig.GetActiveProfile()
}
