package llm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCommandWithArgs(t *testing.T) {
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
