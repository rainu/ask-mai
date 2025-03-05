package controller

import (
	"context"
	"encoding/json"
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
	Id           string
	Role         string
	ContentParts []LLMMessageContentPart
}
type LLMMessageContentPart struct {
	Type    LLMMessageContentPartType
	Content string
	Call    LLMMessageCall
}

type LLMMessageCall struct {
	Id            string
	Function      string
	Arguments     string
	NeedsApproval bool
	Result        *LLMMessageCallResult
}

type LLMMessageCallResult struct {
	Content    string
	Error      string
	DurationMs int64
}

type LLMMessageContentPartType string

const (
	LLMMessageContentPartTypeAttachment LLMMessageContentPartType = "attachment"
	LLMMessageContentPartTypeText       LLMMessageContentPartType = "text"
	LLMMessageContentPartTypeToolCall   LLMMessageContentPartType = "tool"
)

type LLMMessages []LLMMessage

func (m LLMMessages) ToMessageContent(systemPrompt string) ([]llms.MessageContent, error) {
	var result []llms.MessageContent

	for _, msg := range m {
		msgResult := llms.MessageContent{
			Role: llms.ChatMessageType(msg.Role),
		}
		for _, part := range msg.ContentParts {
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
				var msgPart llms.ContentPart = binaryPart

				// special treatment for images (some llms supports image URLs but no binary containing image data)
				if strings.HasPrefix(binaryPart.MIMEType, "image/") {
					msgPart = llms.ImageURLPart(binaryPart.String())
				}
				msgResult.Parts = append(msgResult.Parts, msgPart)
			case LLMMessageContentPartTypeToolCall:
				result = append(result, llms.MessageContent{
					Role: llms.ChatMessageTypeAI,
					Parts: []llms.ContentPart{llms.ToolCall{
						ID:   part.Call.Id,
						Type: string(llms.ChatMessageTypeFunction),
						FunctionCall: &llms.FunctionCall{
							Name:      part.Call.Function,
							Arguments: part.Call.Arguments,
						},
					}},
				})

				if part.Call.Result != nil {
					resultAsJson, err := json.Marshal(map[string]string{
						"output": part.Call.Result.Content,
						"error":  part.Call.Result.Error,
					})
					if err != nil {
						return nil, fmt.Errorf("error serializing tool call result: %w", err)
					}
					msgResult.Parts = append(msgResult.Parts, llms.ToolCallResponse{
						ToolCallID: part.Call.Id,
						Name:       part.Call.Function,
						Content:    string(resultAsJson),
					})
					msgResult.Parts = append(msgResult.Parts)
				}
			case LLMMessageContentPartTypeText:
				fallthrough
			default:
				msgResult.Parts = append(msgResult.Parts, llms.TextPart(part.Content))
			}
		}

		result = append(result, msgResult)
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

	opts := c.appConfig.LLM.CallOptions.AsOptions()
	opts = append(opts, c.appConfig.LLM.Tools.AsOptions()...)
	if c.appConfig.UI.Stream {
		// streaming is enabled
		opts = append(opts, llms.WithStreamingFunc(c.streamingFunc))
	}

	var resp *llms.ContentResponse
	for {
		content, err := args.History.ToMessageContent(c.appConfig.LLM.CallOptions.SystemPrompt)
		if err != nil {
			return "", fmt.Errorf("error converting history to message content: %w", err)
		}

		resp, err = c.aiModel.GenerateContent(c.aiModelCtx, content, opts...)
		if err != nil {
			return "", fmt.Errorf("error creating completion: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no completion choices returned")
		}

		tcMessage, err := c.handleToolCall(resp)
		if err != nil {
			return "", fmt.Errorf("error handling tool call: %w", err)
		}
		if tcMessage != nil {
			args.History = append(args.History, *tcMessage)
			// and continue with the next iteration
		} else {
			break
		}
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
