package controller

import (
	"fmt"
	"github.com/rainu/ask-mai/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"math"
)

func (c *Controller) GetApplicationConfig() config.Config {
	return *c.appConfig
}

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	c.applyInitialWindowConfig()

	runtime.WindowShow(c.ctx)
}

func (c *Controller) applyInitialWindowConfig() {
	screens, err := runtime.ScreenGetAll(c.ctx)
	if err != nil {
		panic(fmt.Errorf("could not get screens: %w", err))
	}

	variables := config.FromScreens(screens)
	_, height := runtime.WindowGetSize(c.ctx)

	value, _ := config.Expression(c.appConfig.UI.Window.InitialWidth).Calculate(variables)
	if int(value) > 0 {
		runtime.WindowSetSize(c.ctx, int(value), height)
	}

	value, _ = config.Expression(c.appConfig.UI.Window.MaxHeight).Calculate(variables)
	if int(value) > 0 {
		runtime.WindowSetMaxSize(c.ctx, math.MaxInt32, int(value))
	}

	posX, _ := config.Expression(c.appConfig.UI.Window.InitialPositionX).Calculate(variables)
	posY, _ := config.Expression(c.appConfig.UI.Window.InitialPositionY).Calculate(variables)

	if int(posX) >= 0 && int(posY) >= 0 {
		runtime.WindowSetPosition(c.ctx, int(posX), int(posY))
	} else {
		runtime.WindowCenter(c.ctx)
	}
}
