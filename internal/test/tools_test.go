package test

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	"github.com/rainu/ask-mai/internal/controller"
	"github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"
	"strings"
	"testing"
	"time"
)

func TestTool_BuiltIn_SystemTime(t *testing.T) {
	cfg, ctrl := initTest(t)

	deactivateAllTools(cfg)
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.SystemTime.Disable = false

	toolCalled := false
	controller.RuntimeEventsEmit = func(ctx context.Context, event string, args ...interface{}) {
		if event == controller.EventNameLLMMessageUpdate {
			msg := args[0].(controller.LLMMessage)
			assert.Equal(t, string(llms.ChatMessageTypeTool), msg.Role)
			require.True(t, strings.HasSuffix(msg.ContentParts[0].Call.Function, system.SystemTimeTool.Name))
			toolCalled = true

			require.NotNil(t, msg.ContentParts[0].Call.Result)
			assert.Contains(t, msg.ContentParts[0].Call.Result.Content, time.Now().String()[:10])
			assert.Empty(t, msg.ContentParts[0].Call.Result.Error)
		}
	}

	res, err := simpleAsk(ctrl, "Wie spät ist es?")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, toolCalled, "Tool should have been called")
}

func TestTool_Custom_Script(t *testing.T) {
	cfg, ctrl := initTest(t)

	deactivateAllTools(cfg)
	cfg.GetActiveProfile().LLM.Tool.Custom[system.SystemTimeTool.Name] = command.FunctionDefinition{
		Description: system.SystemTimeTool.Description,
		CommandExpr: `'2010-08-13T20:15:00Z'`,
	}
	require.NoError(t, cfg.Validate())

	toolCalled := false
	controller.RuntimeEventsEmit = func(ctx context.Context, event string, args ...interface{}) {
		if event == controller.EventNameLLMMessageUpdate {
			msg := args[0].(controller.LLMMessage)
			assert.Equal(t, string(llms.ChatMessageTypeTool), msg.Role)
			require.True(t, strings.HasSuffix(msg.ContentParts[0].Call.Function, system.SystemTimeTool.Name))
			toolCalled = true

			require.NotNil(t, msg.ContentParts[0].Call.Result)
			assert.Contains(t, msg.ContentParts[0].Call.Result.Content, "2010-08-13T20:15:00Z")
			assert.Empty(t, msg.ContentParts[0].Call.Result.Error)
		}
	}

	res, err := simpleAsk(ctrl, "Wie spät ist es?")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, toolCalled, "Tool should have been called")
}

func TestTool_Custom_ScriptWithArguments(t *testing.T) {
	cfg, ctrl := initTest(t)

	deactivateAllTools(cfg)
	cfg.GetActiveProfile().LLM.Tool.Custom[system.SystemTimeTool.Name] = command.FunctionDefinition{
		Description: "Say something to the user.",
		Parameters: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "Message to say to the user.",
				},
			},
			Required: []string{"message"},
		},
		CommandExpr: `
const pArgs = JSON.parse(ctx.args)
'The AI says: ' + pArgs['message']
`,
	}
	require.NoError(t, cfg.Validate())

	toolCalled := false
	controller.RuntimeEventsEmit = func(ctx context.Context, event string, args ...interface{}) {
		if event == controller.EventNameLLMMessageUpdate {
			msg := args[0].(controller.LLMMessage)
			assert.Equal(t, string(llms.ChatMessageTypeTool), msg.Role)
			require.True(t, strings.HasSuffix(msg.ContentParts[0].Call.Function, system.SystemTimeTool.Name))
			toolCalled = true

			require.NotNil(t, msg.ContentParts[0].Call.Result)
			assert.Contains(t, msg.ContentParts[0].Call.Result.Content, `The AI says: Hallo Benutzer!`)
			assert.Empty(t, msg.ContentParts[0].Call.Result.Error)
		}
	}

	res, err := simpleAsk(ctrl, "Sag zu dem Benutzer: Hallo Benutzer!")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, toolCalled, "Tool should have been called")
}

func TestTool_Custom_Command(t *testing.T) {
	cfg, ctrl := initTest(t)

	deactivateAllTools(cfg)
	cfg.GetActiveProfile().LLM.Tool.Custom[system.SystemTimeTool.Name] = command.FunctionDefinition{
		Description: "Say something to the user.",
		Parameters: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "Message to say to the user.",
				},
			},
			Required: []string{"message"},
		},
		Command: `echo "The AI says: $message"`,
	}
	require.NoError(t, cfg.Validate())

	toolCalled := false
	controller.RuntimeEventsEmit = func(ctx context.Context, event string, args ...interface{}) {
		if event == controller.EventNameLLMMessageUpdate {
			msg := args[0].(controller.LLMMessage)
			assert.Equal(t, string(llms.ChatMessageTypeTool), msg.Role)
			require.True(t, strings.HasSuffix(msg.ContentParts[0].Call.Function, system.SystemTimeTool.Name))
			toolCalled = true

			require.NotNil(t, msg.ContentParts[0].Call.Result)
			assert.Contains(t, msg.ContentParts[0].Call.Result.Content, `The AI says: Hallo Benutzer!`)
			assert.Empty(t, msg.ContentParts[0].Call.Result.Error)
		}
	}

	res, err := simpleAsk(ctrl, "Sag zu dem Benutzer: Hallo Benutzer!")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, toolCalled, "Tool should have been called")
}
