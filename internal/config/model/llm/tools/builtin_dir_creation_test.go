package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_DirectoryCreation(t *testing.T) {
	toTest := NewDirectoryCreation()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
