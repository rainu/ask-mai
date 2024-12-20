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

func (m LLMMessages) ToMessageContent(systemPrompt string) []llms.MessageContent {
	result := make([]llms.MessageContent, len(m))
	for i, msg := range m {
		result[i] = llms.MessageContent{
			Parts: []llms.ContentPart{
				llms.TextContent{Text: msg.Content},
			},
			Role: llms.ChatMessageType(msg.Role),
		}
	}

	if systemPrompt != "" {
		msg := llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt)
		result = append([]llms.MessageContent{msg}, result...)
	}

	return result
}

func (c *Controller) LLMAsk(args LLMAskArgs) (result string, err error) {
	defer func() {
		c.aiModelMutex.Write(func() {
			// save the result for later usage (wait)
			c.lastAskResult = llmAskResult{
				Content: result,
				Error:   err,
			}
		})
	}()

	if len(args.History) == 0 {
		return "", fmt.Errorf("empty history provided")
	}
	err = c.LLMInterrupt()
	if err != nil {
		return "", fmt.Errorf("error interrupting previous LLM: %w", err)
	}

	content := args.History.ToMessageContent(c.appConfig.CallOptions.SystemPrompt)

	c.aiModelMutex.Write(func() {
		c.aiModelCtx, c.aiModelCancel = context.WithCancel(context.Background())
	})
	defer func() {
		c.aiModelCancel()
		c.aiModelMutex.Write(func() {
			c.aiModelCtx = nil
			c.aiModelCancel = nil
		})
	}()

	resp, err := c.aiModel.GenerateContent(c.aiModelCtx, content, c.appConfig.CallOptions.AsOptions()...)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	result = resp.Choices[0].Content
	if result != "" {
		c.printer.Print(args.History[len(args.History)-1].Content, result)
	}

	return result, nil
}

func (c *Controller) LLMWait() (result string, err error) {
	var waitChan <-chan struct{}

	c.aiModelMutex.Read(func() {
		if c.aiModelCtx != nil {
			waitChan = c.aiModelCtx.Done()
		}
	})

	// the waiting must be "outside" of the mutex
	if waitChan != nil {
		<-waitChan
	}

	c.aiModelMutex.Read(func() {
		result = c.lastAskResult.Content
		err = c.lastAskResult.Error
	})

	return
}

func (c *Controller) LLMInterrupt() (err error) {
	var cancelFn context.CancelFunc = func() {}

	c.aiModelMutex.Read(func() {
		if c.aiModelCancel != nil {
			cancelFn = c.aiModelCancel
		}
	})

	cancelFn()

	return nil
}
