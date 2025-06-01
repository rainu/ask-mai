package google

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/rainu/ask-mai/internal/sync"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"time"
)

type Google struct {
	client       *googleai.GoogleAI
	clientCtx    context.Context
	clientCancel context.CancelFunc

	cacheTTL     time.Duration
	cacheRefresh time.Duration
	cacheNames   sync.Map[string, string]
}

func New(opts []googleai.Option, toolCacheTTL time.Duration) (common.Model, error) {
	result := &Google{}
	result.clientCtx, result.clientCancel = context.WithCancel(context.Background())

	opts = append(opts, googleai.WithRest())
	if toolCacheTTL > 0 {
		result.cacheRefresh = toolCacheTTL - 30*time.Second
		opts = append(opts, googleai.WithPreSendingHook(result.preSendingHook))
	}

	var err error
	result.client, err = googleai.New(result.clientCtx, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating GoogleAI LLM: %w", err)
	}

	return result, nil
}

func (g *Google) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return g.client.GenerateContent(ctx, messages, options...)
}

func (g *Google) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return g.client.Call(ctx, prompt, options...)
}

func (g *Google) Close() error {
	g.clientCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	g.removeAllCaches(ctx)

	return nil
}
