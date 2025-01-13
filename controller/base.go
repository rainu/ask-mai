package controller

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/config"
	"github.com/rainu/ask-mai/io"
	"github.com/rainu/ask-mai/llms"
	"github.com/rainu/ask-mai/sync"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
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

	screens, err := runtime.ScreenGetAll(ctx)
	if err != nil {
		panic(fmt.Errorf("could not get screens: %w", err))
	}

	err = c.appConfig.ResolveExpressions(config.FromScreens(screens))
	if err != nil {
		panic(fmt.Errorf("could not resolve expressions: %w", err))
	}
}

func (c *Controller) domReady(ctx context.Context) {
}

func (c *Controller) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (c *Controller) Shutdown() {
	c.shutdown(c.ctx)
	os.Exit(0)
}

func (c *Controller) shutdown(ctx context.Context) {
	c.LLMInterrupt()
	c.aiModel.Close()

	for _, target := range c.appConfig.Printer.Targets {
		target.Close()
	}
}
