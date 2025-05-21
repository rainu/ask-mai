package history

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"
	"os"
	"strings"
	"testing"
)

func TestReadWrite(t *testing.T) {
	testHistFile, err := os.CreateTemp(t.TempDir(), "history")
	require.NoError(t, err)

	testWriter := NewWriter(testHistFile.Name())
	testReader := NewReader(testHistFile.Name())

	const testCount = 10
	for i := 1; i <= testCount; i++ {
		we := testWriter.Write(Entry{
			Meta: EntryMeta{
				Version:   1,
				Timestamp: int64(i),
			},
			Content: EntryContent{
				Messages: []Message{
					{
						Role:         string(llms.ChatMessageTypeHuman),
						ContentParts: []MessageContentPart{{Type: "text", Content: "How much time is it?"}},
					},
					{
						Role:         string(llms.ChatMessageTypeAI),
						ContentParts: []MessageContentPart{{Type: "text", Content: "It is 11:09 a.m. on April 2, 2025."}},
					},
				},
			},
		})
		require.NoError(t, we)
	}

	rc, err := testReader.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, testCount, rc)

	entries, err := testReader.GetLast(0, 5)
	assert.NoError(t, err)
	assert.Len(t, entries, 5)
	for i, entry := range entries {
		// we expected the entries to be in reverse order
		assert.Equal(t, entry.Meta.Timestamp, int64(testCount-i))
	}

	entries, err = testReader.GetLast(5, 5)
	assert.NoError(t, err)
	assert.Len(t, entries, 5)
	for i, entry := range entries {
		// we expected the entries to be in reverse order
		assert.Equal(t, entry.Meta.Timestamp, int64(testCount-5-i))
	}

	entries, err = testReader.GetLast(0, testCount*2)
	assert.NoError(t, err)
	assert.Len(t, entries, testCount)
	for i, entry := range entries {
		assert.Equal(t, entry.Meta.Timestamp, int64(testCount-i))
	}

	found := 0
	entries, err = testReader.Search(func(entry Entry) (bool, bool) {
		match := entry.Meta.Timestamp == 7 || entry.Meta.Timestamp == 3 || entry.Meta.Timestamp == 1
		if match {
			found++
		}

		return match, found < 2
	})
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, entries[0].Meta.Timestamp, int64(7))
	assert.Equal(t, entries[1].Meta.Timestamp, int64(3))
}

func TestNewReader_LargeLines(t *testing.T) {
	testHistFile, err := os.CreateTemp(t.TempDir(), "history")
	require.NoError(t, err)

	testWriter := NewWriter(testHistFile.Name())
	testReader := NewReader(testHistFile.Name())

	e := Entry{
		Meta: EntryMeta{},
		Content: EntryContent{
			Messages: []Message{
				{
					Id:           "Id",
					Role:         "Role",
					ContentParts: []MessageContentPart{{Type: "text", Content: strings.Repeat("lorem ipsum", 10000)}},
					CreatedAt:    13,
				},
			},
		},
	}
	require.NoError(t, testWriter.Write(e))

	entries, err := testReader.GetLast(0, 1)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, e, entries[0])
}
