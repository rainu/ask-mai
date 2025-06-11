package command

import (
	"bytes"
	"context"
	cmdchain "github.com/rainu/go-command-chain"
)

type CommandDescriptor struct {
	Command               string            `json:"command"`
	Arguments             []string          `json:"arguments"`
	Environment           map[string]string `json:"env"`
	AdditionalEnvironment map[string]string `json:"additionalEnv"`
	WorkingDirectory      string            `json:"workingDir"`
}

func (c CommandDescriptor) Run(ctx context.Context) ([]byte, error) {
	cmdBuild := cmdchain.Builder().JoinWithContext(ctx, c.Command, c.Arguments...)

	if len(c.Environment) > 0 {
		cmdBuild = cmdBuild.WithEnvironmentMap(toAnyMap(c.Environment))
	}
	if len(c.AdditionalEnvironment) > 0 {
		cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(toAnyMap(c.AdditionalEnvironment))
	}
	if c.WorkingDirectory != "" {
		cmdBuild = cmdBuild.WithWorkingDirectory(c.WorkingDirectory)
	}

	buf := bytes.NewBuffer([]byte{})
	execErr := cmdBuild.WithErrorChecker(cmdchain.IgnoreExitErrors()).Finalize().
		WithOutput(buf).
		WithError(buf).
		Run()

	return buf.Bytes(), execErr
}

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}
