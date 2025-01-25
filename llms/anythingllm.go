package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	ID           string `json:"id"`
	Type         string `json:"type"`
	TextResponse string `json:"textResponse"`
	Close        bool   `json:"close"`
	Error        string `json:"error"`
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

func NewAnythingLLM(baseURL, token, workspace string, deleteThread bool) (Model, error) {
	result := &AnythingLLM{
		client: &http.Client{},

		token:     token,
		baseURL:   baseURL,
		workspace: workspace,

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

	result, err := a.doChat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling anythingllm: %w", err)
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: result},
		},
	}, nil
}

func (a *AnythingLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	err := a.ensureThread(ctx)
	if err != nil {
		return "", err
	}

	return a.doChat(ctx, chatRequest{
		Message:   prompt,
		Mode:      "chat",
		SessionID: a.threadSlug,
	})
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

func (a *AnythingLLM) doChat(ctx context.Context, request chatRequest) (string, error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s/chat", a.baseURL, a.workspace, a.threadSlug)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return result.TextResponse, nil
}

func (a *AnythingLLM) createNewThread(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/new", a.baseURL, a.workspace)

	jsonPayload, err := json.Marshal(threadRequest{
		Name: "ask mAI - " + time.Now().Format("2006-01-02T15:04:05"),
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
