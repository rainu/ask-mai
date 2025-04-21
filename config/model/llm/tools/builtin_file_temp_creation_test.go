package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_FileTempCreation(t *testing.T) {
	toTest := NewFileTempCreation()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
