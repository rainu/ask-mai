package controller

import (
	"context"
	"github.com/rainu/ask-mai/backend"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
)

type Controller struct {
	ctx context.Context

	backendBuilder backend.Builder
	backend        backend.Handle

	appConfig *config.Config
	printer   io.ResponsePrinter
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
}

func (c *Controller) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (c *Controller) shutdown(ctx context.Context) {
	if c.backend != nil {
		c.backend.Close()
	}
	for _, target := range c.appConfig.Printer.Targets {
		target.Close()
	}
}
