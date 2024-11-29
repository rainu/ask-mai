package openai

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/rainu/ask-mai/backend"
)

type OpenAI struct {
	client *openai.Client
	ctx    context.Context
	cancel context.CancelFunc

	history []openai.ChatCompletionMessageParamUnion
}

func NewOpenAI(apiKey, systemPrompt string) (backend.Handle, error) {
	result := &OpenAI{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
	if systemPrompt != "" {
		result.history = append(result.history, openai.AssistantMessage(systemPrompt))
	}

	result.ctx, result.cancel = context.WithCancel(context.Background())
	return result, nil
}

func (o *OpenAI) AskSomething(question string) (string, error) {
	return o.AskSomethingWithContext(o.ctx, question)
}

func (o *OpenAI) AskSomethingWithContext(ctx context.Context, question string) (string, error) {
	o.history = append(o.history, openai.UserMessage(question))

	req := openai.ChatCompletionNewParams{
		Model:    openai.F(openai.ChatModelGPT4oMini),
		Messages: openai.F(o.history),
	}

	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	o.history = append(o.history, openai.AssistantMessage(resp.Choices[0].Message.Content))
	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) Close() error {
	if o.cancel != nil {
		o.cancel()
	}
	return nil
}
