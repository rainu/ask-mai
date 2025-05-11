package mcp

import (
	"fmt"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/metoro-io/mcp-golang/transport/http"
)

type Http struct {
	BaseUrl  string            `yaml:"baseUrl,omitempty" usage:"[http] Base URL"`
	Endpoint string            `yaml:"endpoint,omitempty" usage:"[http] Endpoint"`
	Headers  map[string]string `yaml:"headers,omitempty" usage:"[http] Headers to pass"`
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
