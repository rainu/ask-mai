package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/expression"
	http2 "github.com/rainu/ask-mai/internal/mcp/server/builtin/tools/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCommandExpression_CommandFn(t *testing.T) {
	toTest := Expression(`JSON.stringify(` + expression.VarNameContext + `)`)
	require.NoError(t, toTest.Validate())

	testVars := Variables{
		FunctionDefinition: FunctionDefinition{
			Name:        "test",
			Description: "This is a test",
			Parameters: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"path": map[string]any{
						"type":        "string",
						"description": "The path to the file.",
					},
				},
				Required: []string{"path"},
			},
			Approval: "false",
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
	toTest := Expression(`log("test")`)
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
			expression: `"Echo: " + JSON.parse(` + expression.VarNameContext + `.args).message`,
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
		exec := func(ce Expression) {
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
			ce := Expression(tt.expression)
			require.NoError(t, ce.Validate())

			exec(ce)
		})

		t.Run(fmt.Sprintf("TestCommandExpression_CommandFn_FileReference_%d", i), func(t *testing.T) {
			tmp, err := os.CreateTemp("", "ask-mai-test.*.js")
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(tmp.Name(), []byte(tt.expression), 0666))

			defer os.Remove(tmp.Name())

			ce := Expression(tmp.Name())
			require.NoError(t, ce.Validate())

			exec(ce)
		})
	}
}

func TestCommandExpression_CommandFn_RunCommand(t *testing.T) {
	toTest := Expression(`
const pa = JSON.parse(` + expression.VarNameContext + `.args)
const cmdDescriptor = {
 "command": "echo",
 "arguments": ["Echo:", pa.message]
}

` + expression.FuncNameRun + `(cmdDescriptor)
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{"message": "Hello World"}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)
	assert.Equal(t, `Echo: Hello World`, strings.TrimSpace(string(result)))
}

func TestCommandExpression_CommandFn_RunCommand_WithLimit(t *testing.T) {
	toTest := Expression(`
const pa = JSON.parse(` + expression.VarNameContext + `.args)
const cmdDescriptor = {
 "command": "echo",
 "arguments": ["Echo:", pa.message],
 "output": {
   "firstNBytes": 1,
 }
}

` + expression.FuncNameRun + `(cmdDescriptor)
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{"message": "Hello World"}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)
	assert.Equal(t, "E\n{{ 17 bytes skipped }}", strings.TrimSpace(string(result)))
}

func TestCommandExpression_CommandFn_RunCommand_WithEnv(t *testing.T) {
	toTest := Expression(`
const pa = JSON.parse(` + expression.VarNameContext + `.args)
const cmdDescriptor = {
 "command": "env",
 "env": {"TEST_ENV": "test"},
 "additionalEnv": {"ADDITIONAL_ENV_VAR": "value"},
}

` + expression.FuncNameRun + `(cmdDescriptor)
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
	toTest := Expression(`
` + expression.FuncNameRun + `({
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
	toTest := Expression(`
let result = ""
try {
	result = ` + expression.FuncNameRun + `({
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

func TestCommandExpression_CommandFn_RunFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))

		// Request-Body lesen
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		assert.Equal(t, `{"test":"data"}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()

	toTest := Expression(`
const pa = JSON.parse(` + expression.VarNameContext + `.args)
const callDescriptor = {
 "method": "GET",
 "url": pa.url,
 "header": {
   "X-Test-Header": "test-value"
 },
 "body": JSON.stringify({"test":"data"})
}

const result = ` + expression.FuncNameFetch + `(callDescriptor)

JSON.stringify(result)
`)
	require.NoError(t, toTest.Validate())

	llmArgs := `{"url": "` + server.URL + `"}`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := toTest.CommandFn(FunctionDefinition{})(ctx, llmArgs)
	assert.NoError(t, err)

	var parsedResult http2.CallResult
	err = json.Unmarshal(result, &parsedResult)

	delete(parsedResult.Header, "Date")

	assert.NoError(t, err)
	assert.Equal(t, http2.CallResult{
		StatusCode: 200,
		Status:     "200 OK",
		Header: map[string][]string{
			"Content-Type":   {"application/json"},
			"Content-Length": {"21"},
		},
		Body: `{"message":"Success"}`,
	}, parsedResult)
}
