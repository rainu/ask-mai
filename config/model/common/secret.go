package common

import (
	"bytes"
	"context"
	"fmt"
	cmdchain "github.com/rainu/go-command-chain"
	"time"
)

type Secret struct {
	Plain   string        `config:"" yaml:"plain" usage:" (plain value)"`
	Command SecretCommand `usage:" (command): "`
}

type SecretCommand struct {
	Name   string            `yaml:"name" usage:"name"`
	Args   []string          `yaml:"args" usage:"arguments"`
	Env    map[string]string `yaml:"env" usage:"additional environment variables"`
	NoTrim bool              `yaml:"no-trim" usage:"dont trim spaces from the output"`
}

func (s Secret) Validate() error {
	if s.Plain == "" && s.Command.Name == "" {
		return fmt.Errorf("no plain value or command provided")
	}

	return nil
}

func (s Secret) GetOrPanicWithDefaultTimeout() []byte {
	return s.GetOrPanicWithTimeout(1 * time.Minute)
}

func (s Secret) GetOrPanicWithTimeout(to time.Duration) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	result, err := s.Get(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get secret: %w\n%s", err, result))
	}
	return result
}

func (s Secret) GetOrPanic(ctx context.Context) []byte {
	result, err := s.Get(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get secret: %w\n%s", err, result))
	}
	return result
}

func (s Secret) Get(ctx context.Context) ([]byte, error) {
	if s.Command.Name != "" {
		return s.Command.Get(ctx)
	}

	return []byte(s.Plain), nil
}

func (s SecretCommand) Get(ctx context.Context) ([]byte, error) {
	result := bytes.NewBuffer([]byte{})
	resultErr := bytes.NewBuffer([]byte{})

	c := cmdchain.Builder().JoinWithContext(ctx, s.Name, s.Args...)

	if len(s.Env) > 0 {
		c = c.WithAdditionalEnvironmentMap(toAnyMap(s.Env))
	}

	err := c.Finalize().WithOutput(result).WithError(resultErr).Run()

	if err != nil {
		return []byte(fmt.Sprintf("[OUT] %s\n[ERR} %s", result.String(), resultErr.String())), err
	}

	if !s.NoTrim {
		return bytes.TrimSpace(result.Bytes()), err
	}

	return result.Bytes(), err
}

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}
