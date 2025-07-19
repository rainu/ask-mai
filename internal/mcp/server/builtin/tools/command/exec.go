package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
)

type CommandExecutionArguments struct {
	Command          string            `json:"command"`
	WorkingDirectory string            `json:"working_directory"`
	Environment      map[string]string `json:"environment"`

	DisableOut bool `json:"out"`
	DisableErr bool `json:"err"`
	FirstBytes int  `json:"first"`
	LastBytes  int  `json:"last"`
}

var defaultCommandExecutionArguments = CommandExecutionArguments{
	FirstBytes: 0,
	LastBytes:  1024,
}

var CommandExecutionTool = mcp.NewTool("executeCommand",
	mcp.WithDescription(`Execute a command on the user's system. Concatenations with binary operators like "&&" or "||" are not supported.`),
	mcp.WithString("command",
		mcp.Required(),
		mcp.Description("The shell-like command to execute."),
	),
	mcp.WithString("working_directory",
		mcp.Description("The working directory for the command."),
	),
	mcp.WithObject("environment",
		mcp.Description("Additional environment variables to pass to the command."),
		mcp.AdditionalProperties(map[string]any{"additionalProperties": true}),
	),
	mcp.WithBoolean("out",
		mcp.Description("Whether to disable standard output. Defaults to false."),
		mcp.DefaultBool(false),
	),
	mcp.WithBoolean("err",
		mcp.Description("Whether to disable standard error. Defaults to false."),
		mcp.DefaultBool(false),
	),
	mcp.WithNumber("first",
		mcp.Description("The number of bytes to read from the beginning of the output. -1 means no limit."),
		mcp.DefaultNumber(float64(defaultCommandExecutionArguments.FirstBytes)),
	),
	mcp.WithNumber("last",
		mcp.Description("The number of bytes to read from the end of the output. -1 means no limit."),
		mcp.DefaultNumber(float64(defaultCommandExecutionArguments.LastBytes)),
	),
)

var CommandExecutionToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pArgs := defaultCommandExecutionArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if pArgs.Command == "" {
		return nil, fmt.Errorf("missing parameter: 'command'")
	}

	cmdDesc := CommandDescriptor{
		CommandLine:           pArgs.Command,
		AdditionalEnvironment: pArgs.Environment,
		WorkingDirectory:      pArgs.WorkingDirectory,
		Output: &OutputSettings{
			DisableStdOut: pArgs.DisableOut,
			DisableStdErr: pArgs.DisableErr,
			FirstNBytes:   pArgs.FirstBytes,
			LastNBytes:    pArgs.LastBytes,
		},
	}

	raw, err := cmdDesc.Run(ctx)
	return mcp.NewToolResultText(string(raw)), err
}
