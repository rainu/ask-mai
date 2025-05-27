package test

import (
	"context"
	"github.com/rainu/ask-mai/internal/config"
	"github.com/rainu/ask-mai/internal/config/model"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/command"
	iMcp "github.com/rainu/ask-mai/internal/config/model/llm/tools/mcp"
	"github.com/rainu/ask-mai/internal/controller"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"
	"log/slog"
	"os"
	"testing"
)

func initTest(t *testing.T, args ...string) (*model.Config, *controller.Controller) {
	_, isCI := os.LookupEnv("CI")
	if isCI {
		t.Skip("Skipping test in CI environment")
		return nil, nil
	}

	cfg := config.Parse(args, nil)
	cfg.DebugConfig.LogLevel = "debug"

	if err := cfg.Validate(); err != nil {
		t.Skip("Skipping test, config validation failed:", err)
		return nil, nil
	}
	slog.SetLogLoggerLevel(*cfg.DebugConfig.LogLevelParsed)

	ctrl, err := controller.BuildFromConfig(cfg, "", false)
	require.NoError(t, err)

	controller.RuntimeEventsEmit = func(ctx context.Context, event string, args ...interface{}) {
		slog.Debug("Runtime event emitted", slog.String("event", event), slog.Any("args", args))
	}

	return cfg, ctrl
}

func deactivateAllTools(cfg *model.Config) *model.Config {
	cfg = deactivateBuiltinTools(cfg)
	cfg = deactivateMCPTools(cfg)
	cfg = deactivateCustomTools(cfg)

	return cfg
}

func deactivateBuiltinTools(cfg *model.Config) *model.Config {
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.SystemInfo.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.Environment.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.SystemTime.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.Stats.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.ChangeMode.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.ChangeOwner.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.ChangeTimes.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.FileCreation.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.FileTempCreation.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.FileAppending.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.FileReading.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.FileDeletion.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.DirectoryCreation.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.DirectoryTempCreation.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.DirectoryDeletion.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.CommandExec.Disable = true
	cfg.GetActiveProfile().LLM.Tool.BuiltIns.Http.Disable = true

	return cfg
}

func deactivateMCPTools(cfg *model.Config) *model.Config {
	cfg.GetActiveProfile().LLM.Tool.McpServer = map[string]iMcp.Server{}

	return cfg
}

func deactivateCustomTools(cfg *model.Config) *model.Config {
	cfg.GetActiveProfile().LLM.Tool.Custom = map[string]command.FunctionDefinition{}

	return cfg
}

func simpleAsk(ctrl *controller.Controller, txtContent string) (controller.LLMAskResult, error) {
	return ctrl.LLMAsk(controller.LLMAskArgs{
		History: controller.LLMMessages{{
			Role: string(llms.ChatMessageTypeHuman),
			ContentParts: []controller.LLMMessageContentPart{{
				Type:    controller.LLMMessageContentPartTypeText,
				Content: txtContent,
			}},
		}},
	})
}
