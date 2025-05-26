package mcp

import (
	"fmt"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/rainu/ask-mai/internal/mcp/client"
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

func (s *Server) GetTransport() (transport transport.Interface, err error) {
	if s.Command.Name != "" {
		transport, err = s.Command.GetTransport()
	} else {
		transport, err = s.Http.GetTransport()
	}

	return transport, err
}

func (s *Server) GetTimeouts() (t client.Timeouts) {
	if s.Timeout.Init != nil {
		t.Init = *s.Timeout.Init
	}
	if s.Timeout.List != nil {
		t.List = *s.Timeout.List
	}
	if s.Timeout.Execution != nil {
		t.Execution = *s.Timeout.Execution
	}

	return
}
