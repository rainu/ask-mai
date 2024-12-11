package controller

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
)

type LLMAskArgs struct {
	History LLMMessages
}

type LLMMessage struct {
	Role    string
	Content string
}
type LLMMessages []LLMMessage

func (m LLMMessages) ToMessageContent() []llms.MessageContent {
	result := make([]llms.MessageContent, len(m))
	for i, msg := range m {
		result[i] = llms.MessageContent{
			Parts: []llms.ContentPart{
				llms.TextContent{Text: msg.Content},
			},
			Role: llms.ChatMessageType(msg.Role),
		}
	}
	return result
}

func (c *Controller) LLMAsk(args LLMAskArgs) (string, error) {
	if len(args.History) == 0 {
		return "", fmt.Errorf("empty history provided")
	}
	err := c.LLMInterrupt()
	if err != nil {
		return "", fmt.Errorf("error interrupting previous LLM: %w", err)
	}

	content := args.History.ToMessageContent()

	c.aiModelCtx, c.aiModelCancel = context.WithCancel(context.Background())
	defer func() {
		c.aiModelCancel()
		c.aiModelCtx = nil
		c.aiModelCancel = nil
	}()

	resp, err := c.aiModel.GenerateContent(c.aiModelCtx, content)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	result := resp.Choices[0].Content
	if result != "" {
		c.printer.Print(args.History[len(args.History)-1].Content, result)
	}

	return result, nil
}

func (c *Controller) LLMInterrupt() (err error) {
	if c.aiModelCancel != nil {
		c.aiModelCancel()
	}

	return nil
}
