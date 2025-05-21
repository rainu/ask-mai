package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_FileDeletion(t *testing.T) {
	toTest := NewFileDeletion()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
