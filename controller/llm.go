package controller

import (
	"fmt"
	"github.com/rainu/ask-mai/backend"
)

type LLMAskArgs struct {
	History []backend.Message
}

func (c *Controller) LLMAsk(args LLMAskArgs) (result string, err error) {
	if len(args.History) == 0 {
		return "", fmt.Errorf("empty history provided")
	}

	for _, message := range args.History {
		if message.Role != backend.RoleUser && message.Role != backend.RoleBot {
			return "", fmt.Errorf("unknown role: %s", message.Role)
		}
	}

	var b backend.Handle
	b, err = c.getBackend()
	if err != nil {
		return "", fmt.Errorf("error getting backend: %w", err)
	}
	result, err = b.AskSomething(args.History)

	if result != "" {
		c.printer.Print(args.History[len(args.History)-1].Content, result)
	}

	return
}

func (c *Controller) LLMInterrupt() (err error) {
	if c.backend == nil {
		return nil
	}

	return c.backend.Close()
}
