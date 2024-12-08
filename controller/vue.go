package controller

import (
	"github.com/rainu/ask-mai/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (c *Controller) GetApplicationConfig() config.Config {
	return *c.appConfig
}

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	runtime.WindowShow(c.ctx)
}
