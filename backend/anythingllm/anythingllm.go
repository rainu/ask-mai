package anythingllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/backend"
	"net/http"
	"time"
)

type AnythingLLM struct {
	client *http.Client
	ctx    context.Context
	cancel context.CancelFunc

	token      string
	baseURL    string
	workspace  string
	threadSlug string
}

type chatRequest struct {
	Message   string `json:"message"`
	Mode      string `json:"mode"`
	SessionID string `json:"sessionId,omitempty"`
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

func NewAnythingLLM(baseURL, token, workspace string) (backend.Handle, error) {
	result := &AnythingLLM{
		client: &http.Client{},

		token:     token,
		baseURL:   baseURL,
		workspace: workspace,
	}

	result.ctx, result.cancel = context.WithCancel(context.Background())

	return result, nil
}

func (a *AnythingLLM) AskSomething(question string) (string, error) {
	return a.AskSomethingWithContext(a.ctx, question)
}

func (a *AnythingLLM) AskSomethingWithContext(ctx context.Context, question string) (string, error) {
	if a.threadSlug == "" {
		err := a.createNewThread()
		if err != nil {
			return "", fmt.Errorf("error creating new thread: %w", err)
		}
	}

	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s/chat", a.baseURL, a.workspace, a.threadSlug)
	jsonPayload, err := json.Marshal(chatRequest{
		Message:   question,
		Mode:      "chat",
		SessionID: a.threadSlug,
	})
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

func (a *AnythingLLM) createNewThread() error {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/new", a.baseURL, a.workspace)

	jsonPayload, err := json.Marshal(threadRequest{
		Name: "ask mAI - " + time.Now().Format("2006-01-02T15:04:05"),
	})
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(a.ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
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

func (a *AnythingLLM) deleteThread() error {
	if a.threadSlug == "" {
		return nil
	}

	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s", a.baseURL, a.workspace, a.threadSlug)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodDelete, url, nil)
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
	if a.cancel != nil {
		a.cancel()
	}

	a.ctx, a.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return a.deleteThread()
}
