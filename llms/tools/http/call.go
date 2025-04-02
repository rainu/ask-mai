package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CallDescriptor struct {
	Method     string            `json:"method"`
	Url        string            `json:"url"`
	Header     map[string]string `json:"header"`
	Body       io.Reader         `json:"-"`
	StringBody string            `json:"body"`
}

type CallResult struct {
	StatusCode int                 `json:"status_code"`
	Status     string              `json:"status"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
}

func (c *CallDescriptor) Run(ctx context.Context, client *http.Client) (*CallResult, error) {
	if c.StringBody != "" {
		c.Body = io.NopCloser(strings.NewReader(c.StringBody))
	}

	req, err := http.NewRequestWithContext(ctx, c.Method, c.Url, c.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.Header {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	result := &CallResult{}
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.Header = resp.Header

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	result.Body = string(rawBody)

	return result, err
}
