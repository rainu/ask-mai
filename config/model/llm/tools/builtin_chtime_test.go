package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_ChangeTimes(t *testing.T) {
	toTest := NewChangeTimes()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
