package approval

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestApproval_NeedsApproval_Always(t *testing.T) {
	assert.True(t, Approval(Always).NeedsApproval(context.Background(), ``, nil))
}

func TestApproval_NeedsApproval_Never(t *testing.T) {
	assert.False(t, Approval(Never).NeedsApproval(context.Background(), ``, nil))
}

func TestApproval_NeedsApproval_InvalidJson(t *testing.T) {
	assert.True(t, Approval(expression.VarNameContext+`.raw_args === 'INVALID_JSON'`).NeedsApproval(context.Background(), `INVALID_JSON`, nil))
}

func TestApproval_NeedsApproval_match(t *testing.T) {
	assert.True(t, Approval(expression.VarNameContext+`.definition.name === 'docker'`).NeedsApproval(context.Background(), `{}`, &mcp.Tool{Name: "docker"}))
}

func TestApproval_NeedsApproval_Logging(t *testing.T) {
	origLog := expression.Log
	defer func() {
		expression.Log = origLog
	}()

	called := false
	expression.Log = func(args ...interface{}) {
		assert.Equal(t, "HelloWorld!", args[0])
		called = true
	}

	Approval(expression.FuncNameLog+`('HelloWorld!'); true`).NeedsApproval(context.Background(), ``, nil)

	assert.True(t, called)
}

func TestApproval_NeedsApproval_Run(t *testing.T) {
	tmp, err := os.CreateTemp("", "ask-mai-test-approval")
	require.NoError(t, err)

	tmp.Close()
	os.Remove(tmp.Name())

	_, err = os.Stat(tmp.Name())
	require.Error(t, err, "file should not exist")

	Approval(expression.FuncNameRun+`({
	"command": "touch",
	"arguments": ['`+tmp.Name()+`'],
}); true`).NeedsApproval(context.Background(), ``, nil)

	_, err = os.Stat(tmp.Name())
	assert.NoError(t, err, "file should exist")
	os.Remove(tmp.Name())
}
