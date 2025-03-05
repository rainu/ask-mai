package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rainu/ask-mai/config/llm/tools"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/tmc/langchaingo/llms"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log/slog"
	"sync"
	"time"
)

func (c *Controller) handleToolCall(resp *llms.ContentResponse) (result *LLMMessage, err error) {
	if len(resp.Choices[0].ToolCalls) == 0 {
		return nil, nil
	}

	result = &LLMMessage{
		Id:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Role: string(llms.ChatMessageTypeTool),
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
		if fnDefinition.NeedsApproval {
			c.toolApprovalMutex.Write(func() {
				c.toolApprovalChannel[call.ID] = make(chan bool)
			})
		}

		result.ContentParts = append(result.ContentParts, LLMMessageContentPart{
			Type: LLMMessageContentPartTypeToolCall,
			Call: LLMMessageCall{
				Id:            call.ID,
				NeedsApproval: fnDefinition.NeedsApproval,
				Function:      call.FunctionCall.Name,
				Arguments:     call.FunctionCall.Arguments,
			},
		})
	}
	runtime.EventsEmit(c.ctx, "llm:message:add", result)

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

				for p := range result.ContentParts {
					if result.ContentParts[p].Call.Id == call.ID {
						result.ContentParts[p].Call.Result = &r
						break
					}
				}
				runtime.EventsEmit(c.ctx, "llm:message:update", result)
			})
		}(i)
	}

	wg.Wait()
	return
}

func (c *Controller) callTool(ctx context.Context, call llms.ToolCall, toolDefinition tools.FunctionDefinition) (result LLMMessageCallResult, err error) {
	if toolDefinition.NeedsApproval {
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

	var out []byte
	var execErr error
	if toolDefinition.CommandFn == nil {
		out, execErr, err = c.executeCommand(ctx, toolDefinition, call)
	} else {
		out, execErr = toolDefinition.CommandFn(ctx, call.FunctionCall.Arguments)
	}

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
		result.Error = fmt.Sprintf("Execution error: %s", err.Error())
		err = nil // do not treat execution errors as error - the LLM will receive the error message
	}

	return
}

func (c *Controller) executeCommand(ctx context.Context, toolDefinition tools.FunctionDefinition, call llms.ToolCall) ([]byte, error, error) {
	cmd, args, err := toolDefinition.GetCommandWithArgs(call.FunctionCall.Arguments)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating command for tool '%s': %w", call.FunctionCall.Name, err)
	}

	cmdBuild := cmdchain.Builder().JoinWithContext(ctx, cmd, args...)

	if len(toolDefinition.Environment) > 0 {
		env, err := toolDefinition.GetEnvironment(call.FunctionCall.Arguments)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating environment for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithEnvironmentMap(env)
	}
	if len(toolDefinition.AdditionalEnvironment) > 0 {
		env, err := toolDefinition.GetAdditionalEnvironment(call.FunctionCall.Arguments)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating additional environment for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(env)
	}
	if toolDefinition.WorkingDir != "" {
		wd, err := toolDefinition.GetWorkingDirectory(call.FunctionCall.Arguments)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating working directory for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithWorkingDirectory(wd)
	}

	buf := bytes.NewBuffer([]byte{})
	execErr := cmdBuild.Finalize().
		WithOutput(buf).
		WithError(buf).
		Run()

	return buf.Bytes(), execErr, nil
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
