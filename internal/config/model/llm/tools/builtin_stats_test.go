package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_Stats(t *testing.T) {
	toTest := NewStats()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
