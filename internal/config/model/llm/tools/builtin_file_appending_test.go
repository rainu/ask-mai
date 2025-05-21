package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_FileAppending(t *testing.T) {
	toTest := NewFileAppending()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
