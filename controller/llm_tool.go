package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	for _, call := range resp.Choices[0].ToolCalls {
		fnDefinition, exists := c.appConfig.LLM.Tools.Tools[call.FunctionCall.Name]
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

			r, e := c.callTool(c.aiModelCtx, resp.Choices[0].ToolCalls[i])
			c.aiModelMutex.Write(func() {
				if e != nil {
					err = errors.Join(err, e)
				}

				for p := range result.ContentParts {
					if result.ContentParts[p].Call.Id == resp.Choices[0].ToolCalls[i].ID {
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

func (c *Controller) callTool(ctx context.Context, call llms.ToolCall) (result LLMMessageCallResult, err error) {
	toolDefinition := c.appConfig.LLM.Tools.Tools[call.FunctionCall.Name]

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
				return result, fmt.Errorf("approval for tool '%s' was rejected", call.FunctionCall.Name)
			}
		case <-ctx.Done():
			return result, fmt.Errorf("approval for tool '%s' timed out", call.FunctionCall.Name)
		}
	}

	cmd, args, err := toolDefinition.GetCommandWithArgs(call.FunctionCall.Arguments)
	if err != nil {
		return result, fmt.Errorf("error creating command for tool '%s': %w", call.FunctionCall.Name, err)
	}

	buf := bytes.NewBuffer([]byte{})
	t := time.Now()

	slog.Debug("Start running command.",
		"command", toolDefinition.Command,
		"argument", call.FunctionCall.Arguments,
	)
	cmdBuild := cmdchain.Builder().JoinWithContext(ctx, cmd, args...)

	if len(toolDefinition.Environment) > 0 {
		env, err := toolDefinition.GetEnvironment(call.FunctionCall.Arguments)
		if err != nil {
			return result, fmt.Errorf("error creating environment for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithEnvironmentMap(env)
	}
	if len(toolDefinition.AdditionalEnvironment) > 0 {
		env, err := toolDefinition.GetAdditionalEnvironment(call.FunctionCall.Arguments)
		if err != nil {
			return result, fmt.Errorf("error creating additional environment for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(env)
	}
	if toolDefinition.WorkingDir != "" {
		wd, err := toolDefinition.GetWorkingDirectory(call.FunctionCall.Arguments)
		if err != nil {
			return result, fmt.Errorf("error creating working directory for tool '%s': %w", call.FunctionCall.Name, err)
		}
		cmdBuild = cmdBuild.WithWorkingDirectory(wd)
	}

	err = cmdBuild.Finalize().
		WithGlobalErrorChecker(cmdchain.IgnoreExitErrors()).WithOutput(buf).WithError(buf).
		Run()

	result.Content = buf.String()
	result.DurationMs = time.Since(t).Milliseconds()
	if err != nil {
		result.Error = err.Error()
		err = fmt.Errorf("error calling tool '%s': %w", cmd, err)
	}

	slog.Debug("Command stopped.",
		"command", cmd,
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
