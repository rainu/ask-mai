package controller

import "github.com/rainu/ask-mai/config/model"

func (c *Controller) GetApplicationConfig() model.Config {
	return *c.getConfig()
}

func (c *Controller) SetActiveProfile(profileName string) model.Config {
	c.appConfig.ActiveProfile = profileName
	return c.GetApplicationConfig()
}

func (c *Controller) getConfig() *model.Config {
	return c.appConfig.GetActiveProfile()
}
