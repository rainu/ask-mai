package controller

import (
	"context"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
	"github.com/rainu/ask-mai/llms"
	"github.com/rainu/ask-mai/sync"
)

type Controller struct {
	ctx context.Context

	aiModel       llms.Model
	aiModelCtx    context.Context
	aiModelCancel context.CancelFunc
	aiModelMutex  sync.Mutex
	lastAskResult llmAskResult

	appConfig *config.Config
	printer   io.ResponsePrinter

	vueAppMounted bool
	streamBuffer  []byte
}

type llmAskResult struct {
	Content string
	Error   error
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
