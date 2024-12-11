package controller

import (
	"context"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
	"github.com/rainu/ask-mai/llms"
)

type Controller struct {
	ctx context.Context

	aiModel       llms.Model
	aiModelCtx    context.Context
	aiModelCancel context.CancelFunc

	appConfig *config.Config
	printer   io.ResponsePrinter
}

func (c *Controller) startup(ctx context.Context) {
	c.ctx = ctx
}

func (c *Controller) domReady(ctx context.Context) {
}

func (c *Controller) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (c *Controller) shutdown(ctx context.Context) {
	c.LLMInterrupt()
	c.aiModel.Close()

	for _, target := range c.appConfig.Printer.Targets {
		target.Close()
	}
}
