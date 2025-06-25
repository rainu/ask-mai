package command

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_getOutput(t *testing.T) {
	testCases := []struct {
		name           string
		content        string
		firstNBytes    int
		lastNBytes     int
		expectedOutput string
	}{
		{
			name:           "read all",
			content:        "This is a test output.",
			firstNBytes:    -1,
			lastNBytes:     -1,
			expectedOutput: "This is a test output.",
		},
		{
			name:           "content too short - read all",
			content:        "This is a test output.",
			firstNBytes:    1024,
			lastNBytes:     0,
			expectedOutput: "This is a test output.",
		},
		{
			name:           "content too short - read all",
			content:        "This is a test output.",
			firstNBytes:    0,
			lastNBytes:     1024,
			expectedOutput: "This is a test output.",
		},
		{
			name:           "first n bytes",
			content:        "This is a test output.",
			firstNBytes:    4,
			lastNBytes:     0,
			expectedOutput: "This\n{{ 18 bytes skipped }}",
		},
		{
			name:           "last n bytes",
			content:        "This is a test output.",
			firstNBytes:    0,
			lastNBytes:     7,
			expectedOutput: "{{ 15 bytes skipped }}\noutput.",
		},
		{
			name:           "first and last n bytes",
			content:        "This is a test output.",
			firstNBytes:    4,
			lastNBytes:     7,
			expectedOutput: "This\n{{ 11 bytes skipped }}\noutput.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp(t.TempDir(), tc.name+"_output.txt")
			require.NoError(t, err)
			defer f.Close()

			_, err = f.WriteString(tc.content)
			require.NoError(t, err)

			// reset position
			_, err = f.Seek(0, 0)
			require.NoError(t, err)

			toTest := CommandDescriptor{Output: &OutputSettings{LastNBytes: tc.lastNBytes, FirstNBytes: tc.firstNBytes}}
			result := toTest.getOutput(f)

			require.Equal(t, tc.expectedOutput, string(result))
		})
	}
}
