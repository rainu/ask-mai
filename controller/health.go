package controller

import "github.com/wailsapp/wails/v2/pkg/runtime"

func (c *Controller) TriggerRestart() {
	runtime.EventsEmit(c.ctx, "system:restart")
}

func (c *Controller) Restart(state string) {
	c.lastState = state
	runtime.Hide(c.ctx)
	runtime.Quit(c.ctx)
}

func (c *Controller) GetLastState() string {
	return c.lastState
}
