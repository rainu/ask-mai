package controller

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/config/model"
	"github.com/rainu/ask-mai/expression"
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

	toolApprovalChannel map[string]chan bool
	toolApprovalMutex   sync.Mutex

	appConfig *model.Config
	printer   io.ResponsePrinter

	currentConversation LLMMessages

	vueAppMounted bool
	streamBuffer  []byte
	lastState     string
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

	sv := expression.SetScreens(screens)

	err = c.appConfig.UI.Window.ResolveExpressions(sv)
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
	c.appConfig.Printer.Close()
	c.saveHistory()
	os.Exit(0)
}

func (c *Controller) shutdown(ctx context.Context) {
	c.LLMInterrupt()
	c.aiModel.Close()
}
