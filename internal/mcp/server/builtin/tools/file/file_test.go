package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPermissions(t *testing.T) {
	t.Parallel()

	// Test cases for permission parsing
	testCases := []struct {
		input    string
		expected os.FileMode
	}{
		{"0777", 0777},
		{"0755", 0755},
		{"0644", 0644},
		{"0000", 0000},
		{"", 01111111},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			mode, err := Permission(tc.input).Get(os.FileMode(01111111))
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, mode)
		})
	}
}
