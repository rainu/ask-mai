package llms

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	client *ollama.LLM
}

func NewOllama(opts []ollama.Option) (Model, error) {
	result := &Ollama{}

	var err error
	result.client, err = ollama.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating Ollama LLM: %w", err)
	}

	return result, nil
}

func (o *Ollama) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return o.client.Call(ctx, prompt, options...)
}

func (o *Ollama) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return o.client.GenerateContent(ctx, messages, options...)
}

func (o *Ollama) Close() error {
	return nil
}
