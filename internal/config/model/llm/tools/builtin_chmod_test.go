package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_ChangeMode(t *testing.T) {
	toTest := NewChangeMode()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
