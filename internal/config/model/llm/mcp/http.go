package mcp

import (
	"fmt"
	"github.com/mark3labs/mcp-go/client/transport"
)

type Http struct {
	BaseUrl string            `yaml:"baseUrl,omitempty" usage:"[http] Base URL"`
	Headers map[string]string `yaml:"headers,omitempty" usage:"[http] Headers to pass"`
}

func (h *Http) Validate() error {
	if len(h.BaseUrl) == 0 {
		return fmt.Errorf("baseUrl is required")
	}
	return nil
}

func (h *Http) GetTransport() (transport.Interface, error) {
	var opts []transport.StreamableHTTPCOption

	if len(h.Headers) > 0 {
		opts = append(opts, transport.WithHTTPHeaders(h.Headers))
	}

	return transport.NewStreamableHTTP(h.BaseUrl, opts...)
}
