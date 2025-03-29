package tools

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestFunctionDefinition_GetCommandWithArgs(t *testing.T) {
	tests := []struct {
		command     string
		args        string
		expectCmd   string
		expectArgs  []string
		expectError bool
	}{
		{
			command:    `/usr/bin/echo $msg`,
			args:       `{"msg": "hello world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"hello", "world"},
		},
		{
			command:    `/usr/bin/echo "$msg"`,
			args:       `{"msg": "hello world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"hello world"},
		},
		{
			command:    `/usr/bin/echo "$msg"`,
			args:       `{"msg": 13}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"13"},
		},
		{
			command:    `/usr/bin/echo "$msg"`,
			args:       `{"msg": 13.12}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"13.12"},
		},
		{
			command:    `/usr/bin/echo "$msg"`,
			args:       `{"msg": {"arg1": "hello", "arg2": "world"}}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{`{"arg1":"hello","arg2":"world"}`},
		},
		{
			command:    `/usr/bin/echo "$arg1" "$arg2"`,
			args:       `{"arg1": "hello", "arg2": "world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"hello", "world"},
		},
		{
			command:    `/usr/bin/echo --arg1 "$arg1" --arg2 "$arg2"`,
			args:       `{"arg1": "hello", "arg2": "world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"--arg1", "hello", "--arg2", "world"},
		},
		{
			command:    `/usr/bin/echo --arg1="$arg1" --arg2="$arg2"`,
			args:       `{"arg1": "hello", "arg2": "world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"--arg1=hello", "--arg2=world"},
		},
		{
			command:    `/usr/bin/echo "$@"`,
			args:       `{"msg": {"arg1": "hello", "arg2": "world"}}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{`{"msg": {"arg1": "hello", "arg2": "world"}}`},
		},
		{
			command:    `/usr/bin/echo`,
			args:       `{"msg": {"arg1": "hello", "arg2": "world"}}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{},
		},
		{
			command:    `/usr/bin/echo $doesNotExist`,
			args:       `{"msg": {"arg1": "hello", "arg2": "world"}}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{},
		},
		{
			command:    `$cmd $args`,
			args:       `{"cmd": "/usr/bin/echo", "args": "hello world"}`,
			expectCmd:  "/usr/bin/echo",
			expectArgs: []string{"hello", "world"},
		},
		{
			command:     `/usr/bin/echo $msg`,
			args:        `BROKEN_JSON`,
			expectError: true,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			fn := &FunctionDefinition{Command: tc.command}
			cmd, args, err := fn.GetCommandWithArgs(tc.args)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectCmd, cmd)
			assert.Equal(t, tc.expectArgs, args)
		})
	}
}

func TestFunctionDefinition_GetEnvironment(t *testing.T) {
	tests := []struct {
		env         map[string]string
		args        string
		expectEnv   map[any]any
		expectError bool
	}{
		{
			env: map[string]string{
				"USER": "rainu",
				"ENV1": "$msg",
			},
			args: `{"msg": "hello world"}`,
			expectEnv: map[any]any{
				"USER": "rainu",
				"ENV1": "hello world",
			},
		},
		{
			env: map[string]string{
				"USER": "rainu",
				"ENV1": "$@",
			},
			args: `{"msg": "hello world"}`,
			expectEnv: map[any]any{
				"USER": "rainu",
				"ENV1": `{"msg": "hello world"}`,
			},
		},
		{
			env: map[string]string{
				"USER": "rainu",
				"ENV1": "$@",
			},
			args:        `BROKEN_JSON`,
			expectError: true,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			fn := &FunctionDefinition{Environment: tc.env}
			re, err := fn.GetEnvironment(tc.args)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectEnv, re)
		})

		t.Run(fmt.Sprintf("Additional_%d", i), func(t *testing.T) {
			fn := &FunctionDefinition{AdditionalEnvironment: tc.env}
			re, err := fn.GetAdditionalEnvironment(tc.args)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectEnv, re)
		})
	}
}

func TestFunctionDefinition_GetWorkingDirectory(t *testing.T) {
	tests := []struct {
		workDir     string
		args        string
		expectWD    string
		expectError bool
	}{
		{
			workDir:  "/usr/$user/home",
			args:     `{"user":"rainu"}`,
			expectWD: "/usr/rainu/home",
		},
		{
			workDir:     "/usr/$user/home",
			args:        `BROKEN_JSON`,
			expectError: true,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			fn := &FunctionDefinition{WorkingDir: tc.workDir}
			re, err := fn.GetWorkingDirectory(tc.args)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectWD, re)
		})
	}
}

func TestCommand_BuiltIn(t *testing.T) {
	testFD := FunctionDefinition{
		Command: `echo "$message"`,
	}

	assert.False(t, testFD.IsBuiltIn())
}

func TestCommand_CommandFn(t *testing.T) {
	testFD := FunctionDefinition{
		Command: `echo "$message"`,
	}
	toTest := Command(testFD.Command)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(testFD)(ctx, `{"message": "Hello World"}`)

	assert.NoError(t, err)
	assert.Equal(t, "Hello World", strings.TrimSpace(string(result)))
}
