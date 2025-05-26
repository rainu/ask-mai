package command

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/config/model/llm/tools/approval"
)

const FunctionArgumentNameAll = "@"

type CommandFn func(ctx context.Context, jsonArguments string) ([]byte, error)
type ApprovalFn func(ctx context.Context, jsonArguments string) bool

type FunctionDefinition struct {
	Name        string              `yaml:"-" json:"name" usage:"The name of the function"`
	Description string              `yaml:"description,omitempty" json:"description" usage:"The description of the function"`
	Parameters  mcp.ToolInputSchema `yaml:"parameters,omitempty" json:"parameters" usage:"The parameter definition of the function"`
	Approval    string              `yaml:"approval,omitempty" json:"approval" usage:"Expression to check if user approval is needed before execute this tool"`

	Command               string            `yaml:"command,omitempty,omitempty" json:"command,omitempty" usage:"The command to execute. This is a format string with placeholders for the parameters. Example: /usr/bin/touch $path"`
	CommandExpr           string            `yaml:"commandExpr,omitempty,omitempty" json:"commandExpr,omitempty" usage:"JavaScript expression (or path to JS-file) to execute. See Tool-Help (--help-tool) for more information."`
	Environment           map[string]string `yaml:"env,omitempty,omitempty" json:"env,omitempty" usage:"Environment variables to pass to the command (will overwrite the default environment)"`
	AdditionalEnvironment map[string]string `yaml:"additionalEnv,omitempty,omitempty" json:"additionalEnv,omitempty" usage:"Additional environment variables to pass to the command (will be merged with the default environment)"`
	WorkingDir            string            `yaml:"workingDir,omitempty,omitempty" json:"workingDir,omitempty" usage:"The working directory for the command"`

	//will be filled at runtime (and should not be filled by user in any way)
	CommandFn  CommandFn  `yaml:"-" json:"-"`
	ApprovalFn ApprovalFn `yaml:"-" json:"-"`
}

func (f *FunctionDefinition) NeedApproval(ctx context.Context, jsonArgs string) bool {
	if f.ApprovalFn == nil {
		return approval.Approval(f.Approval).NeedsApproval(ctx, jsonArgs, f)
	}
	return f.ApprovalFn(ctx, jsonArgs)
}
