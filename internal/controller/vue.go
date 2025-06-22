package controller

import (
	"math"
)

const initialConversationIdInitPrompt = "initial-prompt"

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	c.applyInitialWindowConfig()

	RuntimeWindowShow(c.ctx)

	c.vueAppMounted = true

	// after the app is mounted, we have to inform about the initial conversation
	// (otherwise the followup askings, will not contain the initial conversation)
	for _, message := range c.initialConversation {
		if message.Id == initialConversationIdInitPrompt {
			// skip the user prompt, because it will handle differently by the frontend
			continue
		}
		RuntimeEventsEmit(c.ctx, EventNameLLMMessageAdd, message)
	}
	c.initialConversation = nil
}

func (c *Controller) IsAppMounted() bool {
	return c.vueAppMounted
}

func (c *Controller) applyInitialWindowConfig() {
	_, height := RuntimeWindowGetSize(c.ctx)

	initWidth := int(*c.getProfile().UI.Window.InitialWidth.Value)
	if int(initWidth) > 0 {
		RuntimeWindowSetSize(c.ctx, initWidth, height)
	}

	maxHeight := int(*c.getProfile().UI.Window.MaxHeight.Value)
	if maxHeight > 0 {
		RuntimeWindowSetMaxSize(c.ctx, math.MaxInt32, maxHeight)
	}

	posX := int(*c.getProfile().UI.Window.InitialPositionX.Value)
	posY := int(*c.getProfile().UI.Window.InitialPositionY.Value)

	if *c.getProfile().UI.Window.GrowTop {
		posY = posY - height
	}

	if posX >= 0 || posY >= 0 {
		RuntimeWindowSetPosition(c.ctx, posX, posY)
	} else {
		RuntimeWindowCenter(c.ctx)
	}
}
