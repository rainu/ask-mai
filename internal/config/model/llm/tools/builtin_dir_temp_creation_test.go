package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_DirectoryTempCreation(t *testing.T) {
	toTest := NewDirectoryTempCreation()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
