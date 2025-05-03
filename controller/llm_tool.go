package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	internalMcp "github.com/rainu/ask-mai/config/model/llm/mcp"
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
	availableTools := c.getConfig().LLM.Tools.GetTools()
	availableMcpTools, err := c.getConfig().LLM.McpServer.ListTools(c.aiModelCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to list mcp tools: %w", err)
	}

	for callIdx, call := range resp.Choices[0].ToolCalls {
		fnDefinition, isTool := availableTools[call.FunctionCall.Name]
		mcpToolInfo, isMcpTool := availableMcpTools[call.FunctionCall.Name]
		if !isTool && !isMcpTool {
			return nil, fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
		}

		isBuiltIn := false
		needsApproval := false
		mcpToolDescription := ""
		if mcpToolInfo.Description != nil {
			mcpToolDescription = *mcpToolInfo.Description
		}

		//create approval channel for tool calls that need approval
		if isTool {
			isBuiltIn = fnDefinition.IsBuiltIn()
			needsApproval = fnDefinition.NeedApproval(c.aiModelCtx, call.FunctionCall.Arguments)
			if needsApproval {
				c.toolApprovalMutex.Write(func() {
					c.toolApprovalChannel[call.ID] = make(chan bool)
				})
			}
		}

		tcMessage := LLMMessage{
			Id:      fmt.Sprintf("%s-%d", idBase, callIdx),
			Role:    string(llms.ChatMessageTypeTool),
			Created: time.Now().Unix(),
			ContentParts: []LLMMessageContentPart{{
				Type: LLMMessageContentPartTypeToolCall,
				Call: LLMMessageCall{
					Id:                 call.ID,
					NeedsApproval:      needsApproval,
					BuiltIn:            isBuiltIn,
					McpTool:            isMcpTool,
					McpToolName:        mcpToolInfo.Name,
					McpToolDescription: mcpToolDescription,
					Function:           call.FunctionCall.Name,
					Arguments:          call.FunctionCall.Arguments,
				},
			}},
		}
		result = append(result, tcMessage)
		runtime.EventsEmit(c.ctx, "llm:message:add", tcMessage)
	}

	wg := sync.WaitGroup{}

	for i := range resp.Choices[0].ToolCalls {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			call := resp.Choices[0].ToolCalls[i]

			var r LLMMessageCallResult
			var e error

			if tool, ok := availableTools[call.FunctionCall.Name]; ok {
				r, e = c.callTool(c.aiModelCtx, call, tool)
			} else if tool, ok := availableMcpTools[call.FunctionCall.Name]; ok {
				r, e = c.callMcpTool(c.aiModelCtx, call, tool)
			} else {
				e = fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
			}

			c.aiModelMutex.Write(func() {
				if e != nil {
					err = errors.Join(err, e)
				}

				for _, tcMessage := range result {
					for p := range tcMessage.ContentParts {
						if tcMessage.ContentParts[p].Call.Id == call.ID {
							tcMessage.ContentParts[p].Call.Result = &r
							runtime.EventsEmit(c.ctx, "llm:message:update", tcMessage)
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

func (c *Controller) callTool(ctx context.Context, call llms.ToolCall, toolDefinition tools.FunctionDefinition) (result LLMMessageCallResult, err error) {
	if toolDefinition.NeedApproval(ctx, call.FunctionCall.Arguments) {
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

	if execErr != nil {
		result.Error = fmt.Sprintf("Execution error: %s", execErr.Error())
		err = nil // do not treat execution errors as error - the LLM will receive the error message
	}

	slog.Debug("Command stopped.",
		"name", toolDefinition.Name,
		"command", toolDefinition.Command,
		"argument", call.FunctionCall.Arguments,
		"duration", result.DurationMs,
		"error", result.Error,
	)

	return
}

func (c *Controller) callMcpTool(ctx context.Context, call llms.ToolCall, toolDefinition internalMcp.Tool) (result LLMMessageCallResult, err error) {
	slog.Debug("Start calling mcp tool.", "name", toolDefinition.Name)

	transport := toolDefinition.Transport.GetTransport()
	defer transport.Close()

	mcpClient := mcp.NewClient(transport)
	_, err = mcpClient.Initialize(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to initialize mcp client: %w", err)
	}

	var args any
	err = json.Unmarshal([]byte(call.FunctionCall.Arguments), &args)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal tool call arguments: %w", err)
	}

	t := time.Now()
	resp, callErr := mcpClient.CallTool(ctx, toolDefinition.Name, &args)

	result.DurationMs = time.Since(t).Milliseconds()
	content, _ := json.Marshal(resp)
	result.Content = string(content)

	if callErr != nil {
		result.Error = fmt.Sprintf("Execution error: %s", callErr.Error())
		err = nil // do not treat execution errors as error - the LLM will receive the error message
	}

	slog.Debug("MCP tool stopped.",
		"name", toolDefinition.Name,
		"argument", call.FunctionCall.Arguments,
		"duration", result.DurationMs,
		"error", result.Error,
	)

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
