package tools

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltin_FileReading(t *testing.T) {
	toTest := NewFileReading()

	assert.False(t, toTest.AsFunctionDefinition().NeedApproval(context.Background(), ``))
}
