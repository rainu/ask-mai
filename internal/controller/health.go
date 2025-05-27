package controller

func (c *Controller) TriggerRestart() {
	RuntimeEventsEmit(c.ctx, "system:restart")
}

func (c *Controller) Restart(state string) {
	c.lastState = state
	RuntimeHide(c.ctx)
	RuntimeQuit(c.ctx)
}

func (c *Controller) GetLastState() string {
	return c.lastState
}
