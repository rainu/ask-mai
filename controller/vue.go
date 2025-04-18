package controller

import (
	"github.com/rainu/ask-mai/config/model"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"math"
)

func (c *Controller) GetApplicationConfig() model.Config {
	return *c.appConfig
}

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	c.applyInitialWindowConfig()

	runtime.WindowShow(c.ctx)

	c.vueAppMounted = true
}

func (c *Controller) IsAppMounted() bool {
	return c.vueAppMounted
}

func (c *Controller) applyInitialWindowConfig() {
	_, height := runtime.WindowGetSize(c.ctx)

	initWidth := int(c.appConfig.UI.Window.InitialWidth.Value)
	if int(initWidth) > 0 {
		runtime.WindowSetSize(c.ctx, initWidth, height)
	}

	maxHeight := int(c.appConfig.UI.Window.MaxHeight.Value)
	if maxHeight > 0 {
		runtime.WindowSetMaxSize(c.ctx, math.MaxInt32, maxHeight)
	}

	posX := int(c.appConfig.UI.Window.InitialPositionX.Value)
	posY := int(c.appConfig.UI.Window.InitialPositionY.Value)

	if c.appConfig.UI.Window.GrowTop {
		posY = posY - height
	}

	if posX >= 0 || posY >= 0 {
		runtime.WindowSetPosition(c.ctx, posX, posY)
	} else {
		runtime.WindowCenter(c.ctx)
	}
}
