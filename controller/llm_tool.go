package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/rainu/ask-mai/config/model/llm/tools"
	"github.com/tmc/langchaingo/llms"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log/slog"
	"sync"
	"time"
)

func (c *Controller) handleToolCall(resp *llms.ContentResponse) (result LLMMessages, err error) {
	if len(resp.Choices[0].ToolCalls) == 0 {
		return nil, nil
	}

	tcMessage := &LLMMessage{
		Id:      fmt.Sprintf("%d", time.Now().UnixNano()),
		Role:    string(llms.ChatMessageTypeTool),
		Created: time.Now().Unix(),
	}

	if resp.Choices[0].Content != "" {
		txtMessage := LLMMessage{
			Id:   tcMessage.Id + "-0",
			Role: string(llms.ChatMessageTypeAI),
			ContentParts: []LLMMessageContentPart{{
				Type:    LLMMessageContentPartTypeText,
				Content: resp.Choices[0].Content,
			}},
			Created: time.Now().Unix(),
		}

		runtime.EventsEmit(c.ctx, "llm:message:add", txtMessage)
		result = append(result, txtMessage)
	}

	c.toolApprovalMutex.Write(func() {
		c.toolApprovalChannel = map[string]chan bool{}
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
	availableTools := c.appConfig.LLM.Tools.GetTools()
	for _, call := range resp.Choices[0].ToolCalls {
		fnDefinition, exists := availableTools[call.FunctionCall.Name]
		if !exists {
			return nil, fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
		}

		//create approval channel for tool calls that need approval
		needsApproval := fnDefinition.CheckApproval(c.aiModelCtx, call.FunctionCall.Arguments)
		if needsApproval {
			c.toolApprovalMutex.Write(func() {
				c.toolApprovalChannel[call.ID] = make(chan bool)
			})
		}

		tcMessage.ContentParts = append(tcMessage.ContentParts, LLMMessageContentPart{
			Type: LLMMessageContentPartTypeToolCall,
			Call: LLMMessageCall{
				Id:            call.ID,
				NeedsApproval: needsApproval,
				BuiltIn:       fnDefinition.IsBuiltIn(),
				Function:      call.FunctionCall.Name,
				Arguments:     call.FunctionCall.Arguments,
			},
		})
	}
	runtime.EventsEmit(c.ctx, "llm:message:add", tcMessage)

	wg := sync.WaitGroup{}

	for i := range resp.Choices[0].ToolCalls {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			call := resp.Choices[0].ToolCalls[i]
			r, e := c.callTool(c.aiModelCtx, call, availableTools[call.FunctionCall.Name])
			c.aiModelMutex.Write(func() {
				if e != nil {
					err = errors.Join(err, e)
				}

				for p := range tcMessage.ContentParts {
					if tcMessage.ContentParts[p].Call.Id == call.ID {
						tcMessage.ContentParts[p].Call.Result = &r
						break
					}
				}
				runtime.EventsEmit(c.ctx, "llm:message:update", tcMessage)
			})
		}(i)
	}

	wg.Wait()

	result = append(result, *tcMessage)
	return
}

func (c *Controller) callTool(ctx context.Context, call llms.ToolCall, toolDefinition tools.FunctionDefinition) (result LLMMessageCallResult, err error) {
	if toolDefinition.CheckApproval(ctx, call.FunctionCall.Arguments) {
		// wait for user's approval (see llmApplyToolCallApproval())
		var approvalChan chan bool
		c.toolApprovalMutex.Read(func() {
			approvalChan = c.toolApprovalChannel[call.ID]
		})

		if approvalChan == nil {
			return result, fmt.Errorf("approval channel for tool '%s' not found", call.FunctionCall.Name)
		}

		// wait for approval
		select {
		case approved := <-approvalChan:
			slog.Debug("Approval received for tool.", "tool", call.FunctionCall.Name, "approved", approved)

			if !approved {
				result.Error = "The user rejected the tool call!"
				return result, nil
			}
		case <-ctx.Done():
			return result, fmt.Errorf("approval for tool '%s' timed out", call.FunctionCall.Name)
		}
	}

	slog.Debug("Start running command.",
		"name", toolDefinition.Name,
		"command", toolDefinition.Command,
		"argument", call.FunctionCall.Arguments,
	)
	t := time.Now()

	out, execErr := toolDefinition.CommandFn(ctx, call.FunctionCall.Arguments)

	result.DurationMs = time.Since(t).Milliseconds()
	result.Content = string(out)

	slog.Debug("Command stopped.",
		"name", toolDefinition.Name,
		"command", toolDefinition.Command,
		"argument", call.FunctionCall.Arguments,
		"duration", result.DurationMs,
		"error", result.Error,
	)

	if err != nil {
		return
	}

	if execErr != nil {
		result.Error = fmt.Sprintf("Execution error: %s", execErr.Error())
		err = nil // do not treat execution errors as error - the LLM will receive the error message
	}

	return
}

func (c *Controller) LLMApproveToolCall(callId string) {
	c.llmApplyToolCallApproval(callId, true)
}

func (c *Controller) LLMRejectToolCall(callId string) {
	c.llmApplyToolCallApproval(callId, false)
}

func (c *Controller) llmApplyToolCallApproval(callId string, approve bool) {
	c.toolApprovalMutex.Read(func() {
		if c.toolApprovalChannel[callId] != nil {
			c.toolApprovalChannel[callId] <- approve
		}
	})
}
