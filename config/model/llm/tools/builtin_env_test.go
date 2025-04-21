package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_Environment(t *testing.T) {
	toTest := NewEnvironment()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
