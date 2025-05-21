package openai

import (
	"context"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAI struct {
	client *openai.LLM
}

func New(opts []openai.Option) (common.Model, error) {
	result := &OpenAI{}

	var err error

	// prevent streaming of function calls
	opts = append(opts, openai.WithStreamingChunkFilter(func(chunkMeta openai.StreamingChunkMetaData) bool {
		if chunkMeta.IsFunctionCall {
			return false
		}
		if chunkMeta.IsToolCall {
			return false
		}
		return true
	}))

	result.client, err = openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating OpenAI LLM: %w", err)
	}

	return result, nil
}

func (o *OpenAI) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *OpenAI) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(ctx, messages, options...)
}

func (o *OpenAI) Close() error {
	return nil
}
