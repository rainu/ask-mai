package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/http"
	"github.com/rainu/go-yacl"
)

type Http struct {
	BaseUrl  string            `yaml:"baseUrl,omitempty" usage:"Base URL for the command"`
	Endpoint string            `yaml:"endpoint,omitempty" usage:"Endpoint of the HTTP server"`
	Headers  map[string]string `yaml:"headers,omitempty" usage:"Headers to pass to the HTTP server"`
	Approval string            `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute a tool"`
	Exclude  []string          `yaml:"exclude,omitempty" usage:"List of tools that should be excluded"`

	Timeout Timeout `yaml:"timeout,omitempty"`
}

func (h *Http) Validate() error {
	if len(h.BaseUrl) == 0 && len(h.Endpoint) == 0 {
		return fmt.Errorf("either baseUrl or endpoint must be set")
	}
	return nil
}

func (h *Http) GetTransport() transport.Transport {
	t := http.NewHTTPClientTransport(h.Endpoint)

	if len(h.BaseUrl) > 0 {
		t = t.WithBaseURL(h.BaseUrl)
	}

	for k, v := range h.Headers {
		t = t.WithHeader(k, v)
	}

	return t
}

func (h *Http) ListTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := h.GetTransport()
	defer t.Close()

	if yacl.D(h.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *h.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listTools(ctx, t, h.Exclude)
}

func (h *Http) ListAllTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := h.GetTransport()
	defer t.Close()

	if yacl.D(h.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *h.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listAllTools(ctx, t)
}
