package controller

import (
	"fmt"
	"github.com/rainu/ask-mai/backend"
)

type LLMAskArgs struct {
	Content string
}

func (c *Controller) LLMAsk(args LLMAskArgs) (result string, err error) {
	var b backend.Handle
	b, err = c.getBackend()
	if err != nil {
		return "", fmt.Errorf("error getting backend: %w", err)
	}
	result, err = b.AskSomething(args.Content)

	if result != "" {
		c.printer.Print(args.Content, result)
	}

	return
}

func (c *Controller) LLMInterrupt() (err error) {
	if c.backend == nil {
		return nil
	}

	return c.backend.Close()
}
