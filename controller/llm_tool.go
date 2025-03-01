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

	//validate tool calls
	for _, call := range resp.Choices[0].ToolCalls {
		_, exists := c.appConfig.LLM.Tools.Tools[call.FunctionCall.Name]
		if !exists {
			return nil, fmt.Errorf("unknown tool: %s", call.FunctionCall.Name)
		}

		result.ContentParts = append(result.ContentParts, LLMMessageContentPart{
			Type: LLMMessageContentPartTypeToolCall,
			Call: LLMMessageCall{
				Id:        call.ID,
				Function:  call.FunctionCall.Name,
				Arguments: call.FunctionCall.Arguments,
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
	cmd := c.appConfig.LLM.Tools.Tools[call.FunctionCall.Name].Command

	buf := bytes.NewBuffer([]byte{})
	t := time.Now()

	slog.Debug("Start running command.",
		"command", cmd,
		"argument", call.FunctionCall.Arguments,
	)
	err = cmdchain.Builder().
		JoinWithContext(ctx, cmd, call.FunctionCall.Arguments).
		Finalize().
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
