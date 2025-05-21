package mcp

import (
	"context"
	"time"
)

type Timeouts struct {
	Init      time.Duration
	List      time.Duration
	Execution time.Duration
}

func (t Timeouts) InitContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return t.wrapContext(ctx, t.Init)
}

func (t Timeouts) ListContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return t.wrapContext(ctx, t.List)
}

func (t Timeouts) ExecutionContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return t.wrapContext(ctx, t.Execution)
}

func (t Timeouts) wrapContext(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	if duration > 0 {
		return context.WithTimeout(ctx, duration)
	}
	return ctx, func() {}
}
