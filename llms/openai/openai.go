package openai

import (
	"context"
	"fmt"
	illms "github.com/rainu/ask-mai/llms"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAI struct {
	client        *openai.LLM
	systemMessage llms.MessageContent
}

func NewOpenAI(apiKey, systemPrompt string) (illms.Model, error) {
	result := &OpenAI{
		systemMessage: llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
	}

	var err error
	result.client, err = openai.New(openai.WithToken(apiKey), openai.WithModel("gpt-4o-mini"))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenAI LLM: %w", err)
	}

	return result, nil
}

func (o *OpenAI) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *OpenAI) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(
		ctx,
		append([]llms.MessageContent{o.systemMessage}, messages...),
		options...,
	)
}

func (o *OpenAI) Close() error {
	return nil
}
