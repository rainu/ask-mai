package controller

import "github.com/wailsapp/wails/v2/pkg/runtime"

func (c *Controller) GetInitialPrompt() string {
	return c.initialPrompt
}

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	runtime.WindowShow(c.ctx)
}
