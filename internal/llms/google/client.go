package google

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type Google struct {
	client       *googleai.GoogleAI
	clientCtx    context.Context
	clientCancel context.CancelFunc
}

func New(opts []googleai.Option) (common.Model, error) {
	result := &Google{}
	result.clientCtx, result.clientCancel = context.WithCancel(context.Background())

	opts = append(opts, googleai.WithRest())

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
	return nil
}
