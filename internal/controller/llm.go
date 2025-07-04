package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/go-yacl"
	"github.com/tmc/langchaingo/llms"
	"os"
	"strings"
	"time"
)

type LLMAskArgs struct {
	History LLMMessages
}

type LLMMessage struct {
	Id           string
	Role         string
	ContentParts []LLMMessageContentPart
	Created      int64
}
type LLMMessageContentPart struct {
	Type    LLMMessageContentPartType
	Content string
	Call    LLMMessageCall
}

type LLMMessageCall struct {
	Id        string
	Function  string
	Arguments string
	Meta      LLMMessageCallMeta
	Result    *LLMMessageCallResult
}

type LLMMessageCallMeta struct {
	BuiltIn bool
	Custom  bool
	Mcp     bool

	NeedsApproval bool

	ToolName        string
	ToolDescription string
}

type LLMMessageCallResult struct {
	Content    string
	Error      string
	DurationMs int64

	//only for wails to generate TypeScript types
	V mcp.CallToolResult
	W mcp.TextContent
	X mcp.ImageContent
	Y mcp.AudioContent
	Z mcp.EmbeddedResource
}

type LLMMessageContentPartType string

const (
	LLMMessageContentPartTypeAttachment LLMMessageContentPartType = "attachment"
	LLMMessageContentPartTypeText       LLMMessageContentPartType = "text"
	LLMMessageContentPartTypeToolCall   LLMMessageContentPartType = "tool"
)

type LLMMessages []LLMMessage

func (m LLMMessages) ToMessageContent() ([]llms.MessageContent, error) {
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
				if part.Content != "" {
					msgResult.Parts = append(msgResult.Parts, llms.TextPart(part.Content))
				}
			}
		}

		if len(msgResult.Parts) > 0 {
			result = append(result, msgResult)
		}
	}

	return result, nil
}

type llmAskResponse struct {
	Result LLMAskResult
	Error  error
}

type LLMAskResult struct {
	Content     string
	Consumption common.ConsumptionSummary
}

func (c *Controller) LLMAsk(args LLMAskArgs) (result LLMAskResult, err error) {
	defer func() {
		c.currentConversation = args.History
		if err == nil {
			c.currentConversation = append(c.currentConversation, LLMMessage{
				Role: string(llms.ChatMessageTypeAI),
				ContentParts: []LLMMessageContentPart{{
					Type:    LLMMessageContentPartTypeText,
					Content: result.Content,
				}},
				Created: time.Now().Unix(),
			})
		}
	}()
	defer func() {
		c.aiModelMutex.Write(func() {
			// save the result for later usage (wait)
			c.lastAskResponse = llmAskResponse{
				Result: result,
				Error:  err,
			}
		})
	}()

	if len(args.History) == 0 {
		return LLMAskResult{}, fmt.Errorf("empty history provided")
	}
	err = c.LLMInterrupt()
	if err != nil {
		return LLMAskResult{}, fmt.Errorf("error interrupting previous LLM: %w", err)
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

	opts, err := c.getProfile().LLM.AsOptions(c.aiModelCtx, c.aiModel)
	if err != nil {
		return LLMAskResult{}, fmt.Errorf("error creating options: %w", err)
	}

	if yacl.D(c.getProfile().UI.Stream) {
		// streaming is enabled
		opts = append(opts, llms.WithStreamingFunc(c.streamingFunc))
	}

	consumption := c.aiModel.ConsumptionOf(nil)
	var resp *llms.ContentResponse
	for {
		content, err := args.History.ToMessageContent()
		if err != nil {
			return LLMAskResult{}, fmt.Errorf("error converting history to message content: %w", err)
		}

		resp, err = c.aiModel.GenerateContent(c.aiModelCtx, content, opts...)
		if err != nil {
			return LLMAskResult{}, fmt.Errorf("error creating completion: %w", err)
		}

		if len(resp.Choices) == 0 {
			return LLMAskResult{}, fmt.Errorf("no completion choices returned")
		}

		consumption.Add(c.aiModel.ConsumptionOf(resp))
		RuntimeEventsEmit(c.ctx, EventNameLLMConsumptionUpdate, consumption.Summary())

		tcMessage, err := c.handleToolCall(resp)
		if err != nil {
			return LLMAskResult{}, fmt.Errorf("error handling tool call: %w", err)
		}
		if len(tcMessage) > 0 {
			args.History = append(args.History, tcMessage...)
			// and continue with the next iteration
		} else {
			break
		}
	}

	result.Content = resp.Choices[0].Content
	result.Consumption = consumption.Summary()
	if result.Content != "" {
		question := ""
		for i := len(args.History) - 1; i >= 0; i-- {
			if args.History[i].Role == string(llms.ChatMessageTypeHuman) {
				question = args.History[i].ContentParts[0].Content
				break
			}
		}

		c.printer.Print(question, result.Content)
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
		RuntimeEventsEmit(c.ctx, "llm:stream:chunk", string(c.streamBuffer))
		c.streamBuffer = nil
	}

	RuntimeEventsEmit(c.ctx, "llm:stream:chunk", string(chunk))
	return nil
}

func (c *Controller) LLMWait() (result LLMAskResult, err error) {
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
		result = c.lastAskResponse.Result
		err = c.lastAskResponse.Error
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
