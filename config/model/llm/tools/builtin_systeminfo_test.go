package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_SystemInfo(t *testing.T) {
	toTest := NewSystemInfo()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
