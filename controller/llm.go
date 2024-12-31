package controller

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/tmc/langchaingo/llms"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"strings"
)

type LLMAskArgs struct {
	History LLMMessages
}

type LLMMessage struct {
	Role         string
	ContentParts []LLMMessageContentPart
}
type LLMMessageContentPart struct {
	Type    LLMMessageContentPartType
	Content string
}
type LLMMessageContentPartType string

const (
	LLMMessageContentPartTypeAttachment LLMMessageContentPartType = "attachment"
	LLMMessageContentPartTypeText       LLMMessageContentPartType = "text"
)

type LLMMessages []LLMMessage

func (m LLMMessages) ToMessageContent(systemPrompt string) ([]llms.MessageContent, error) {
	result := make([]llms.MessageContent, len(m))
	for i, msg := range m {
		result[i] = llms.MessageContent{
			Role:  llms.ChatMessageType(msg.Role),
			Parts: make([]llms.ContentPart, len(msg.ContentParts)),
		}
		for j, part := range msg.ContentParts {
			switch part.Type {
			case LLMMessageContentPartTypeAttachment:
				path := part.Content // in case of attachment, the content is the file path
				mime, err := mimetype.DetectFile(path)
				if err != nil {
					return nil, fmt.Errorf("error detecting mimetype: %w", err)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("error reading attachment: %w", err)
				}
				binaryPart := llms.BinaryPart(mime.String(), data)
				result[i].Parts[j] = binaryPart

				// special treatment for images (some llms supports image URLs but no binary containing image data)
				if strings.HasPrefix(binaryPart.MIMEType, "image/") {
					result[i].Parts[j] = llms.ImageURLPart(binaryPart.String())
				}
			case LLMMessageContentPartTypeText:
				fallthrough
			default:
				result[i].Parts[j] = llms.TextPart(part.Content)
			}
		}
	}

	if systemPrompt != "" {
		msg := llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt)
		result = append([]llms.MessageContent{msg}, result...)
	}

	return result, nil
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

	content, err := args.History.ToMessageContent(c.appConfig.CallOptions.SystemPrompt)
	if err != nil {
		return "", fmt.Errorf("error converting history to message content: %w", err)
	}

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

	opts := c.appConfig.CallOptions.AsOptions()
	if c.appConfig.UI.Stream {
		// streaming is enabled
		opts = append(opts, llms.WithStreamingFunc(c.streamingFunc))
	}

	resp, err := c.aiModel.GenerateContent(c.aiModelCtx, content, opts...)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	result = resp.Choices[0].Content
	if result != "" {
		c.printer.Print(args.History[len(args.History)-1].ContentParts[0].Content, result)
	}

	return result, nil
}

func (c *Controller) streamingFunc(ctx context.Context, chunk []byte) error {
	if !c.vueAppMounted {
		c.streamBuffer = append(c.streamBuffer, chunk...)
		return nil
	}
	if len(c.streamBuffer) > 0 {
		// emit the buffered chunk
		runtime.EventsEmit(c.ctx, "llm:stream:chunk", string(c.streamBuffer))
		c.streamBuffer = nil
	}

	runtime.EventsEmit(c.ctx, "llm:stream:chunk", string(chunk))
	return nil
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
