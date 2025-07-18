package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools"
	"github.com/rainu/ask-mai/internal/mcp/client"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const (
	EventNameLLMMessageAdd        = "llm:message:add"
	EventNameLLMMessageUpdate     = "llm:message:update"
	EventNameLLMConsumptionUpdate = "llm:consumption:update"
)

func (c *Controller) handleToolCall(resp *llms.ContentResponse) (result LLMMessages, err error) {
	if len(resp.Choices[0].ToolCalls) == 0 {
		return nil, nil
	}

	idBase := fmt.Sprintf("%d", time.Now().UnixNano())

	if resp.Choices[0].Content != "" {
		txtMessage := LLMMessage{
			Id:   idBase + "-t",
			Role: string(llms.ChatMessageTypeAI),
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: resp.Choices[0].Content,
			}},
			Created: time.Now().Unix(),
		}

		RuntimeEventsEmit(c.ctx, EventNameLLMMessageAdd, txtMessage)
		result = append(result, txtMessage)
	}

	c.toolApprovalMutex.Write(func() {
		c.toolApprovalChannel = map[string]chan approval{}
	})

	// close approval channel to prevent memory leaks
	defer func() {
		c.toolApprovalMutex.Read(func() {
			for _, c := range c.toolApprovalChannel {
				close(c)
			}
		})
	}()

	//validate tool calls
	availableTools, err := c.getProfile().LLM.Tool.GetTools(c.aiModelCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	for callIdx, call := range resp.Choices[0].ToolCalls {
		toolInfo, ok := availableTools[call.FunctionCall.Name]
		if !ok {
			return nil, fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
		}

		needsApproval := toolInfo.NeedApproval(c.aiModelCtx, call.FunctionCall.Arguments)

		//create approval channel for tool calls that need approval
		if needsApproval {
			c.toolApprovalMutex.Write(func() {
				c.toolApprovalChannel[call.ID] = make(chan approval)
			})
		}

		tcMessage := LLMMessage{
			Id:      fmt.Sprintf("%s-%d", idBase, callIdx),
			Role:    string(llms.ChatMessageTypeTool),
			Created: time.Now().Unix(),
			ContentParts: []LLMMessageContentPart{{
				Type: LLMMessageContentPartTypeToolCall,
				Call: LLMMessageCall{
					Id:        call.ID,
					Function:  call.FunctionCall.Name,
					Arguments: call.FunctionCall.Arguments,

					Meta: LLMMessageCallMeta{
						BuiltIn:         strings.HasPrefix(call.FunctionCall.Name, tools.ServerNameBuiltin),
						Custom:          strings.HasPrefix(call.FunctionCall.Name, tools.ServerNameCustom),
						Mcp:             !strings.HasPrefix(call.FunctionCall.Name, tools.ServerNameBuiltin) && !strings.HasPrefix(call.FunctionCall.Name, tools.ServerNameCustom),
						NeedsApproval:   needsApproval,
						ToolName:        toolInfo.Name,
						ToolDescription: toolInfo.Description,
					},
				},
			}},
		}
		result = append(result, tcMessage)
		RuntimeEventsEmit(c.ctx, EventNameLLMMessageAdd, tcMessage)
	}

	wg := sync.WaitGroup{}

	for i := range resp.Choices[0].ToolCalls {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			call := resp.Choices[0].ToolCalls[i]

			var r LLMMessageCallResult
			var e error

			aErr := c.waitForApproval(c.aiModelCtx, call)
			if aErr != nil {
				r.Error = aErr.Error()
			} else {
				if tool, ok := availableTools[call.FunctionCall.Name]; ok {
					r = c.callTool(c.aiModelCtx, call, tool)
				} else {
					e = fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
				}
			}

			c.aiModelMutex.Write(func() {
				if e != nil {
					err = errors.Join(err, e)
				}

				for _, tcMessage := range result {
					for p := range tcMessage.ContentParts {
						if tcMessage.ContentParts[p].Call.Id == call.ID {
							tcMessage.ContentParts[p].Call.Result = &r
							RuntimeEventsEmit(c.ctx, EventNameLLMMessageUpdate, tcMessage)
							return
						}
					}
				}
			})
		}(i)
	}

	wg.Wait()

	return
}

func (c *Controller) callTool(ctx context.Context, call llms.ToolCall, toolDefinition tools.Tool) LLMMessageCallResult {
	return callTool(ctx, toolDefinition.Transporter, toolDefinition.Name, call.FunctionCall.Arguments)
}

func callTool(ctx context.Context, tp client.Transporter, name string, args string) (result LLMMessageCallResult) {
	slog.Debug("Start calling tool.", "name", name)

	t := time.Now()
	resp, callErr := client.CallTool(ctx, tp, name, args)

	result.DurationMs = time.Since(t).Milliseconds()
	content, _ := json.Marshal(resp)
	result.Content = string(content)

	if callErr != nil {
		result.Error = fmt.Sprintf("Execution error: %s", callErr.Error())
	}

	slog.Debug("Tool stopped.",
		"name", name,
		"argument", args,
		"duration", result.DurationMs,
		"error", result.Error,
	)

	return
}
