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

	systemPrompt string
}

func NewOpenAI(apiKey, systemPrompt string) (backend.Handle, error) {
	result := &OpenAI{
		client:       openai.NewClient(option.WithAPIKey(apiKey)),
		systemPrompt: systemPrompt,
	}

	return result, nil
}

func (o *OpenAI) AskSomething(chat []backend.Message) (string, error) {
	o.Close()

	o.ctx, o.cancel = context.WithCancel(context.Background())
	return o.AskSomethingWithContext(o.ctx, chat)
}

func (o *OpenAI) AskSomethingWithContext(ctx context.Context, chat []backend.Message) (string, error) {
	if len(chat) == 0 {
		return "", fmt.Errorf("empty history provided")
	}

	history := make([]openai.ChatCompletionMessageParamUnion, len(chat)+1)

	history[0] = openai.SystemMessage(o.systemPrompt)
	for i, msg := range chat {
		switch msg.Role {
		case backend.RoleUser:
			history[i+1] = openai.UserMessage(msg.Content)
		case backend.RoleBot:
			history[i+1] = openai.AssistantMessage(msg.Content)
		default:
			return "", fmt.Errorf("unknown role: %s", msg.Role)
		}
	}

	req := openai.ChatCompletionNewParams{
		Model:    openai.F(openai.ChatModelGPT4oMini),
		Messages: openai.F(history),
	}

	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) Close() error {
	if o.cancel != nil {
		o.cancel()
	}
	return nil
}
