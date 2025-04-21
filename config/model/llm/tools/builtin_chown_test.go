package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_ChangeOwner(t *testing.T) {
	toTest := NewChangeOwner()

	assert.True(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
