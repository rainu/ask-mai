package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/config/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCommandExpression_CommandFn(t *testing.T) {
	toTest := CommandExpression(`JSON.stringify(v)`)
	require.NoError(t, toTest.Validate())

	testVars := CommandVariables{
		FunctionDefinition: FunctionDefinition{
			Name:        "test",
			Description: "This is a test",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file.",
					},
				},
				"required": []string{"path"},
			},
			NeedsApproval: false,
			Environment: map[string]string{
				"TEST_ENV": "test",
			},
			AdditionalEnvironment: map[string]string{
				"ADDITIONAL_ENV_VAR": "value",
			},
			WorkingDir:  "/home/test",
			Command:     "EMPTY",
			CommandExpr: string(toTest),
		},
		Arguments: `{"path": "/tmp/"}`,
	}
	varsAsJson, err := json.Marshal(testVars)
	require.NoError(t, err)

	testFn := toTest.CommandFn(testVars.FunctionDefinition)

	result, err := testFn(context.Background(), testVars.Arguments)

	assert.NoError(t, err)
	assert.JSONEq(t, string(varsAsJson), string(result), "Parameter seems not to be passed correctly")
}

func TestCommandExpression_CommandFn_InternalLog(t *testing.T) {
	toTest := CommandExpression(`log("test")`)
	require.NoError(t, toTest.Validate())

	origLog := expression.Log
	defer func() {
		expression.Log = origLog
	}()

	var logCalledArgs []any
	expression.Log = func(args ...interface{}) {
		logCalledArgs = args
	}

	_, err := toTest.CommandFn(FunctionDefinition{})(context.Background(), "{}")
	assert.NoError(t, err)
	assert.Equal(t, logCalledArgs, []any{"test"})
}

func TestCommandExpression_CommandFn_Functionality(t *testing.T) {
	tests := []struct {
		expression string
		args       string
		expected   string
		assertion  func(t *testing.T, result []byte)
	}{
		{
			expression: `"test"`,
			expected:   "test",
		},
		{
			expression: `"Echo: " + JSON.parse(v.args).message`,
			args:       `{"message": "Hello World"}`,
			expected:   `Echo: Hello World`,
		},
		{
			expression: `
let r = ""
for (let i = 0; i < 10; i++) { 
	r += " " + i 
}
r.trim()`,
			expected: "0 1 2 3 4 5 6 7 8 9",
		},
		{
			expression: `new Date().getTime()`,
			assertion: func(t *testing.T, result []byte) {
				assert.Regexp(t, `^[0-9]{13}$`, string(result))
			},
		},
	}

	for i, tt := range tests {
		exec := func(ce CommandExpression) {
			jsonArg := tt.args
			if jsonArg == "" {
				jsonArg = "{}"
			}

			result, err := ce.CommandFn(FunctionDefinition{})(context.Background(), jsonArg)
			assert.NoError(t, err)

			if tt.assertion != nil {
				tt.assertion(t, result)
			} else {
				assert.Equal(t, tt.expected, string(result))
			}
		}

		t.Run(fmt.Sprintf("TestCommandExpression_CommandFn_%d", i), func(t *testing.T) {
			ce := CommandExpression(tt.expression)
			require.NoError(t, ce.Validate())

			exec(ce)
		})

		t.Run(fmt.Sprintf("TestCommandExpression_CommandFn_FileReference_%d", i), func(t *testing.T) {
			tmp, err := os.CreateTemp("", "ask-mai-test.*.js")
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(tmp.Name(), []byte(tt.expression), 0666))

			defer os.Remove(tmp.Name())

			ce := CommandExpression(tmp.Name())
			require.NoError(t, ce.Validate())

			exec(ce)
		})
	}
}

func TestCommandExpression_CommandFn_RunCommand(t *testing.T) {
	toTest := CommandExpression(`
const pa = JSON.parse(v.args)
const cmdDescriptor = {
 "command": "echo",
 "arguments": ["Echo:", pa.message]
}

` + FuncNameRun + `(cmdDescriptor)
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{"message": "Hello World"}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)
	assert.Equal(t, `Echo: Hello World`, strings.TrimSpace(string(result)))
}

func TestCommandExpression_CommandFn_RunCommand_WithEnv(t *testing.T) {
	toTest := CommandExpression(`
const pa = JSON.parse(v.args)
const cmdDescriptor = {
 "command": "env",
 "env": {"TEST_ENV": "test"},
 "additionalEnv": {"ADDITIONAL_ENV_VAR": "value"},
}

` + FuncNameRun + `(cmdDescriptor)
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{"message": "Hello World"}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)
	assert.Contains(t, strings.TrimSpace(string(result)), `TEST_ENV=test`)
	assert.Contains(t, strings.TrimSpace(string(result)), `ADDITIONAL_ENV_VAR=value`)
}

func TestCommandExpression_CommandFn_RunCommand_WithError(t *testing.T) {
	toTest := CommandExpression(`
` + FuncNameRun + `({
 "command": "__DoesNotExistOnAnySystem__"
})
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "executable file not found in $PATH")
}

func TestCommandExpression_CommandFn_RunCommand_CatchError(t *testing.T) {
	toTest := CommandExpression(`
let result = ""
try {
	result = ` + FuncNameRun + `({
		"command": "__DoesNotExistOnAnySystem__"
	})
} catch (e) {
	result = "Error: " + e
}

result
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)
	assert.Equal(t, `Error: failed to start command: exec: "__DoesNotExistOnAnySystem__": executable file not found in $PATH`, strings.TrimSpace(string(result)))
}
