package mcp

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"log/slog"
	"sync"
)

type Transporter interface {
	GetTransport() (transport.Interface, error)
	GetTimeouts() Timeouts
}

type contextAndCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var clientPool = map[string]client.MCPClient{}
var clientTransportContext = map[string]*contextAndCancel{}
var clientPoolMutex = &sync.Mutex{}

func GetClient(ctx context.Context, tp Transporter) (c client.MCPClient, err error) {
	key := fmt.Sprintf("%p", tp)

	clientPoolMutex.Lock()
	defer clientPoolMutex.Unlock()

	c, ok := clientPool[key]
	if ok {
		return c, nil
	}

	t, err := tp.GetTransport()
	if err != nil {
		return nil, fmt.Errorf("failed to get transport: %w", err)
	}

	clientTransportContext[key] = &contextAndCancel{}
	clientTransportContext[key].ctx, clientTransportContext[key].cancel = context.WithCancel(context.Background())
	go func() {
		<-clientTransportContext[key].ctx.Done()
		slog.Debug("MCP transport context was closed!", "key", key)

		clientPoolMutex.Lock()
		c.Close()
		delete(clientTransportContext, key)
		clientPoolMutex.Unlock()
	}()

	err = t.Start(clientTransportContext[key].ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transport: %w", err)
	}

	c = client.NewClient(t)

	initCtx, cancel := tp.GetTimeouts().InitContext(ctx)
	defer cancel()

	_, err = c.Initialize(initCtx, mcp.InitializeRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mcp client: %w", err)
	}

	clientPool[key] = c
	return c, nil
}

func Close() {
	clientPoolMutex.Lock()
	defer clientPoolMutex.Unlock()

	for _, c := range clientPool {
		c.Close()
	}
}
