package anythingllm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/common"
	"github.com/tmc/langchaingo/llms"
	"net/http"
	"strings"
	"time"
)

type AnythingLLM struct {
	client *http.Client

	token     string
	baseURL   string
	workspace string

	threadSlug   string
	threadName   string
	threadDelete bool
}

type chatRequest struct {
	Message     string           `json:"message"`
	Attachments []chatAttachment `json:"attachments,omitempty"`
	Mode        string           `json:"mode"`
	SessionID   string           `json:"sessionId,omitempty"`
}

type chatAttachment struct {
	Name    string `json:"name"`
	Mime    string `json:"mime"`
	Content string `json:"contentString"`
}

type chatResponse struct {
	ID           string      `json:"id"`
	ChatId       int         `json:"chatId"`
	Type         string      `json:"type"`
	TextResponse string      `json:"textResponse"`
	Close        bool        `json:"close"`
	Error        string      `json:"error"`
	Metrics      chatMetrics `json:"metrics"`
}

const streamChatDataEventPrefix = "data: "

type streamChatResponse struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`

	TextResponse string `json:"textResponse"`

	Metrics chatMetrics `json:"metrics"`
	ChatId  int         `json:"chatId"`

	Close bool `json:"close"`
	Error bool `json:"error"`
}

type chatMetrics struct {
	PromptTokens      int     `json:"prompt_tokens"`
	CompletionTokens  int     `json:"completion_tokens"`
	TotalTokens       int     `json:"total_tokens"`
	OutputTokenPerSec float64 `json:"outputTps"`
	Duration          float64 `json:"duration"`
}

type threadRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type threadResponse struct {
	Thread struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		UserID      int    `json:"user_id"`
		WorkspaceID int    `json:"workspace_id"`
	} `json:"thread"`
	Message *string `json:"message"`
}

func New(baseURL, token, workspace, threadName string, deleteThread bool) (common.Model, error) {
	result := &AnythingLLM{
		client: &http.Client{},

		token:     token,
		baseURL:   baseURL,
		workspace: workspace,

		threadName:   threadName,
		threadDelete: deleteThread,
	}

	return result, nil
}

func (a *AnythingLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("empty messages provided")
	}

	err := a.ensureThread(ctx)
	if err != nil {
		return nil, err
	}

	req := chatRequest{Mode: "chat", SessionID: a.threadSlug}
	for _, part := range messages[len(messages)-1].Parts {
		switch p := part.(type) {
		case llms.TextContent:
			req.Message += p.Text
		case llms.ImageURLContent:
			attachment := chatAttachment{
				Name:    fmt.Sprintf("image_%d", len(req.Attachments)),
				Content: p.URL,
			}
			if strings.HasPrefix(p.URL, "data:") {
				attachment.Mime = strings.Split(p.URL[len("data:"):], ";")[0]
			}

			req.Attachments = append(req.Attachments, attachment)
		}
	}
	opts := getOpts(options)

	var result *chatResponse
	if opts.StreamingFunc == nil {
		result, err = a.doChat(ctx, req)
	} else {
		result, err = a.doStreamChat(ctx, req, opts.StreamingFunc)
	}
	if err != nil {
		return nil, fmt.Errorf("error calling anythingllm: %w", err)
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: result.TextResponse,
				GenerationInfo: map[string]any{
					"id":                  result.ID,
					"chatId":              result.ChatId,
					generalInfoKeyMetrics: result.Metrics,
				},
			},
		},
	}, nil
}

func (a *AnythingLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	err := a.ensureThread(ctx)
	if err != nil {
		return "", err
	}

	opts := getOpts(options)
	req := chatRequest{
		Message:   prompt,
		Mode:      "chat",
		SessionID: a.threadSlug,
	}

	var result *chatResponse
	if opts.StreamingFunc == nil {
		result, err = a.doChat(ctx, req)
	} else {
		result, err = a.doStreamChat(ctx, req, opts.StreamingFunc)
	}

	if err != nil {
		return "", fmt.Errorf("error calling anythingllm: %w", err)
	}

	return result.TextResponse, nil
}

func getOpts(options []llms.CallOption) llms.CallOptions {
	opts := llms.CallOptions{}
	for _, o := range options {
		o(&opts)
	}
	return opts
}

func (a *AnythingLLM) ensureThread(ctx context.Context) error {
	if a.threadSlug == "" {
		err := a.createNewThread(ctx)
		if err != nil {
			return fmt.Errorf("error creating new thread: %w", err)
		}
	}
	return nil
}

func (a *AnythingLLM) doChat(ctx context.Context, request chatRequest) (*chatResponse, error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s/chat", a.baseURL, a.workspace, a.threadSlug)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func (a *AnythingLLM) doStreamChat(ctx context.Context, request chatRequest, streamFn func(ctx context.Context, chunk []byte) error) (*chatResponse, error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s/stream-chat", a.baseURL, a.workspace, a.threadSlug)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result := &chatResponse{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, streamChatDataEventPrefix) {
			continue
		}

		line = line[len(streamChatDataEventPrefix):]

		var parsed streamChatResponse
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		if parsed.Error {
			return nil, fmt.Errorf("error in response: %s", line)
		}

		switch parsed.Type {
		case "textResponseChunk":
			se := streamFn(ctx, []byte(parsed.TextResponse))
			if se != nil {
				return nil, fmt.Errorf("error in streaming function: %w", se)
			}
			result.TextResponse += parsed.TextResponse
		case "finalizeResponseStream":
			result.ChatId = parsed.ChatId
			result.Metrics = parsed.Metrics
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return result, nil
}

func (a *AnythingLLM) createNewThread(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/new", a.baseURL, a.workspace)

	jsonPayload, err := json.Marshal(threadRequest{
		Name: a.threadName,
	})
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result threadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	a.threadSlug = result.Thread.Slug

	return nil
}

func (a *AnythingLLM) deleteThread(ctx context.Context) error {
	if a.threadSlug == "" {
		return nil
	}

	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s", a.baseURL, a.workspace, a.threadSlug)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	a.threadSlug = ""

	return nil
}

func (a *AnythingLLM) Close() error {
	if !a.threadDelete {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.deleteThread(ctx)
}
