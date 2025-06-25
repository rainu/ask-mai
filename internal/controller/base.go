package controller

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/expression"
	"github.com/rainu/ask-mai/internal/io"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/ask-mai/internal/sync"
	"os"
)

type Controller struct {
	ctx context.Context

	aiModel         common.Model
	aiModelCtx      context.Context
	aiModelCancel   context.CancelFunc
	aiModelMutex    sync.Mutex
	lastAskResponse llmAskResponse

	toolApprovalChannel map[string]chan approval
	toolApprovalMutex   sync.Mutex

	appConfig *model.Config
	printer   io.ResponsePrinter

	initialConversation LLMMessages
	currentConversation LLMMessages

	vueAppMounted bool
	streamBuffer  []byte
	lastState     string
}

func (c *Controller) startup(ctx context.Context) {
	c.ctx = ctx

	screens, err := RuntimeScreenGetAll(ctx)
	if err != nil {
		panic(fmt.Errorf("could not get screens: %w", err))
	}

	sv := expression.SetScreens(screens)

	err = c.appConfig.MainProfile.UI.Window.ResolveExpressions(sv)
	if err != nil {
		panic(fmt.Errorf("could not resolve expressions: %w", err))
	}
	for profile, config := range c.appConfig.Profiles {
		err = config.UI.Window.ResolveExpressions(sv)
		if err != nil {
			panic(fmt.Errorf("could not resolve expressions from profile '%s': %w", profile, err))
		}
	}
}

func (c *Controller) domReady(ctx context.Context) {
}

func (c *Controller) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (c *Controller) Shutdown() {
	c.shutdown(c.ctx)

	c.appConfig.MainProfile.Printer.Close()
	for _, config := range c.appConfig.Profiles {
		config.Printer.Close()
	}

	c.saveHistory()
	os.Exit(0)
}

func (c *Controller) shutdown(ctx context.Context) {
	c.LLMInterrupt()
	c.aiModel.Close()
}
