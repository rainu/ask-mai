package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_FileCreation(t *testing.T) {
	toTest := NewFileCreation()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
