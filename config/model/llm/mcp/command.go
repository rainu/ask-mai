package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/rainu/go-yacl"
	"io"
	"log/slog"
)

type Command struct {
	Name                  string            `yaml:"name,omitempty" usage:"Name of the command to execute"`
	Arguments             []string          `yaml:"args,omitempty" usage:"Arguments to pass to the command"`
	Environment           map[string]string `yaml:"env,omitempty" usage:"Environment variables to pass to the command"`
	AdditionalEnvironment map[string]string `yaml:"additionalEnv,omitempty" usage:"Additional environment variables to pass to the command"`
	WorkingDirectory      string            `yaml:"workingDir,omitempty" usage:"Working directory for the command"`
	Approval              string            `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute a tool"`
	Exclude               []string          `yaml:"exclude,omitempty" usage:"List of tools that should be excluded"`

	Timeout Timeout `yaml:"timeout,omitempty"`
}

func (c *Command) Validate() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("command name is required")
	}
	return nil
}

type transportProxy struct {
	transport.Transport
	closeHandler func()
}

func (t *transportProxy) Close() error {
	t.closeHandler()
	return t.Transport.Close()
}

func (c *Command) GetTransport() transport.Transport {
	rIn, wIn := io.Pipe()
	rOut, wOut := io.Pipe()

	t := &transportProxy{Transport: stdio.NewStdioServerTransportWithIO(rIn, wOut)}
	t.closeHandler = func() {
		wIn.Close()
		wOut.Close()
	}

	go func() {
		defer t.Close()

		cmd := cmdchain.Builder().WithInput(rOut).Join(c.Name, c.Arguments...)

		if len(c.Environment) > 0 {
			cmd = cmd.WithEnvironmentMap(toAnyMap(c.Environment))
		}
		if len(c.AdditionalEnvironment) > 0 {
			cmd = cmd.WithAdditionalEnvironmentMap(toAnyMap(c.AdditionalEnvironment))
		}
		if c.WorkingDirectory != "" {
			cmd = cmd.WithWorkingDirectory(c.WorkingDirectory)
		}

		err := cmd.Finalize().WithOutput(wIn).Run()
		if err != nil {
			slog.Warn("Error while executing mcp-command", "command", c.Name, "error", err)
		}
	}()

	return t
}

func (c *Command) ListTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := c.GetTransport()
	defer t.Close()

	if yacl.D(c.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *c.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listTools(ctx, t, c.Exclude)
}

func (c *Command) ListAllTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := c.GetTransport()
	defer t.Close()

	if yacl.D(c.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *c.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listAllTools(ctx, t)
}

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}
