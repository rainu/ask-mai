package mcp

import (
	"fmt"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	cmdchain "github.com/rainu/go-command-chain"
	"io"
	"log/slog"
)

type Command struct {
	Name                  string            `yaml:"command,omitempty" usage:"[Command] Name of the command"`
	Arguments             []string          `yaml:"args,omitempty" usage:"[Command] Arguments to pass"`
	Environment           map[string]string `yaml:"env,omitempty" usage:"[Command] Environment variables to pass"`
	AdditionalEnvironment map[string]string `yaml:"aenv,omitempty" usage:"[Command] Additional environment variables to pass"`
	WorkingDirectory      string            `yaml:"dir,omitempty" usage:"[Command] Working directory"`
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

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}
