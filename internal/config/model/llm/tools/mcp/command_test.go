package mcp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		command Command
		wantErr bool
	}{
		{
			name: "valid command",
			command: Command{
				Name: "echo",
			},
			wantErr: false,
		},
		{
			name: "empty command name",
			command: Command{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.command.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Command.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommand_GetTransport(t *testing.T) {
	toTest := Command{
		Name: "docker",
		Arguments: []string{
			"run", "--rm", "-i", "-e", "GITHUB_PERSONAL_ACCESS_TOKEN=github_", "ghcr.io/github/github-mcp-server",
		},
	}

	result, err := toTest.GetTransport()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
