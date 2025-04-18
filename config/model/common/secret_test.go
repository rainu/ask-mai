package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecret_Get(t *testing.T) {
	tests := []struct {
		conf     Secret
		expected []byte
	}{
		{
			conf: Secret{
				Plain: "secret",
			},
			expected: []byte("secret"),
		},
		{
			conf: Secret{
				Command: SecretCommand{
					Name: "echo",
					Args: []string{"secret"},
				},
			},
			expected: []byte("secret"),
		},
		{
			conf: Secret{
				Command: SecretCommand{
					Name:   "echo",
					Args:   []string{"-n", "secret"},
					NoTrim: true,
				},
			},
			expected: []byte("secret"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestSecret_Get_%d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.conf.GetOrPanicWithDefaultTimeout())
		})
	}
}
