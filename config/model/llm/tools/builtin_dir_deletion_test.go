package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_DirectoryDeletion(t *testing.T) {
	toTest := NewDirectoryDeletion()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
