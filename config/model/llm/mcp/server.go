package mcp

import (
	"context"
	"fmt"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
	"github.com/rainu/go-yacl"
)

type Server struct {
	Http    `yaml:",inline"`
	Command `yaml:",inline"`

	Approval string   `yaml:"approval,omitempty" usage:"Expression to check if user approval is needed before execute a tool"`
	Exclude  []string `yaml:"exclude,omitempty" usage:"List of tools that should be excluded"`

	Timeout Timeout `yaml:"timeout,omitempty"`
}

func (s *Server) Validate() error {
	if s.Command.Name != "" {
		if ve := s.Command.Validate(); ve != nil {
			return ve
		}
	} else if s.Http.BaseUrl != "" {
		if ve := s.Http.Validate(); ve != nil {
			return ve
		}
	} else {
		return fmt.Errorf("either command-name or http-url must be set")
	}

	return nil
}

func (s *Server) GetTransport() transport.Transport {
	if s.Command.Name != "" {
		return s.Command.GetTransport()
	}
	return s.Http.GetTransport()
}

func (s *Server) ListTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := s.GetTransport()
	defer t.Close()

	if yacl.D(s.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *s.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listTools(ctx, t, s.Exclude)
}

func (s *Server) ListAllTools(ctx context.Context) ([]mcp.ToolRetType, error) {
	t := s.GetTransport()
	defer t.Close()

	if yacl.D(s.Timeout.List) > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, *s.Timeout.List)
		defer cancel()

		ctx = ctxWithTimeout
	}

	return listAllTools(ctx, t)
}
