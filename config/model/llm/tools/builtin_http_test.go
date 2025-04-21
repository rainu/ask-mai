package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_Http(t *testing.T) {
	toTest := NewHttp()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
