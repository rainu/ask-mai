package controller

import (
	"fmt"
	"math"
	"time"
)

// AppMounted is called when the frontend application is mounted (decided by the frontend itself)
func (c *Controller) AppMounted() {
	c.applyInitialWindowConfig()

	RuntimeWindowShow(c.ctx)

	c.vueAppMounted = true

	// after the app is mounted, we have to inform about the initial messages
	//(otherwise the followup askings, will not contain the initial messages)
	for i, message := range c.getProfile().LLM.CallOptions.Prompt.InitMessages {
		RuntimeEventsEmit(c.ctx, EventNameLLMMessageAdd, LLMMessage{
			Id:      fmt.Sprintf("initial-%d", i),
			Role:    string(message.Role),
			Created: time.Now().Unix(),
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: message.Content,
			}},
		})
	}
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
