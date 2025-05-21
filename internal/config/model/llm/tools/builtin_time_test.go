package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_SystemTime(t *testing.T) {
	toTest := NewSystemTime()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
