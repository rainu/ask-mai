package controller

import (
	"context"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/io"
)

type Controller struct {
	ctx context.Context

	backendBuilder backend.Builder
	backend        backend.Handle

	initialPrompt string
	printer       io.ResponsePrinter
}

func New(bb backend.Builder, printer io.ResponsePrinter, prompt string) *Controller {
	return &Controller{
		backendBuilder: bb,
		printer:        printer,
		initialPrompt:  prompt,
	}
}

func (c *Controller) getBackend() (backend.Handle, error) {
	if c.backend == nil {
		var err error
		c.backend, err = c.backendBuilder.Build()
		if err != nil {
			return nil, err
		}
	}

	return c.backend, nil
}

func (c *Controller) startup(ctx context.Context) {
	c.ctx = ctx
}

func (c *Controller) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

func (c *Controller) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (c *Controller) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}
