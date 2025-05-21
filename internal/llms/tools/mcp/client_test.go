package mcp

import (
	"github.com/mark3labs/mcp-go/client/transport"
	cmdchain "github.com/rainu/go-command-chain"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type testTransport struct{}

func (t testTransport) GetTransport() (transport.Interface, error) {
	return transport.NewStdio("docker", nil, "run", "--rm", "-i", "-e", "BRAVE_API_KEY=BS...", "mcp/brave-search"), nil
}

func (t testTransport) GetTimeouts() Timeouts {
	return Timeouts{}
}

func TestGetClient_ReEstablishing(t *testing.T) {
	_, isCI := os.LookupEnv("CI")
	if isCI {
		t.Skip("Skipping test in CI environment")
		return
	}

	_, _, err := cmdchain.Builder().Join("docker", "-v").Finalize().RunAndGet()
	if err != nil {
		t.Skipf("Docker is not available: %v", err)
		return
	}

	tt := &testTransport{}

	c, err := GetClient(t.Context(), tt)
	assert.NoError(t, err)

	assert.NoError(t, c.Ping(t.Context()))

	c2, err := GetClient(t.Context(), tt)
	assert.NoError(t, err)
	assert.Same(t, c, c2)

	assert.NoError(t, c2.Ping(t.Context()))

	// now interrupt the underlying process and check if we get a new client
	for _, cc := range clientTransportContext {
		cc.cancel()
	}

	// we have to wait a little so the context-checker can close the client
	time.Sleep(100 * time.Millisecond)

	c3, err := GetClient(t.Context(), tt)
	assert.NoError(t, err)
	assert.NotSame(t, c, c3)

	assert.NoError(t, c3.Ping(t.Context()))
}
