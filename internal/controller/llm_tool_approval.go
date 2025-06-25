package controller

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
)

type approval struct {
	ToolCallId string
	Approved   bool
	Message    string
}

func (c *Controller) waitForApproval(ctx context.Context, call llms.ToolCall) error {
	// wait for user's approval (see llmApplyToolCallApproval())
	var approvalChan chan approval
	c.toolApprovalMutex.Read(func() {
		approvalChan = c.toolApprovalChannel[call.ID]
	})

	if approvalChan != nil {
		defer func() {
			c.toolApprovalMutex.Write(func() {
				delete(c.toolApprovalChannel, call.ID)
			})
		}()

		// wait for approval
		select {
		case a := <-approvalChan:
			slog.Debug("Approval received for tool.",
				"tool", call.FunctionCall.Name,
				"approved", a.Approved,
				"message", a.Message,
			)

			if !a.Approved {
				errMsg := "The user rejected the tool call!"
				if a.Message != "" {
					errMsg = a.Message
				}
				return fmt.Errorf(errMsg)
			}
		case <-ctx.Done():
			return fmt.Errorf("Approval for tool '%s' timed out!", call.FunctionCall.Name)
		}
	}

	// no approval needed
	return nil
}

func (c *Controller) LLMApproveToolCall(callId, message string) {
	c.llmApplyToolCallApproval(callId, true, message)
}

func (c *Controller) LLMRejectToolCall(callId, message string) {
	c.llmApplyToolCallApproval(callId, false, message)
}

func (c *Controller) llmApplyToolCallApproval(callId string, approve bool, message string) {
	c.toolApprovalMutex.Read(func() {
		if c.toolApprovalChannel[callId] != nil {
			c.toolApprovalChannel[callId] <- approval{
				ToolCallId: callId,
				Approved:   approve,
				Message:    message,
			}
			close(c.toolApprovalChannel[callId])
		}
	})
}
