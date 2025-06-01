package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map[string, int]{}

	assert.Equal(t, 0, m.Len())

	// Test Write and Read
	m.Write("key1", 1)
	val, ok := m.Read("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
	assert.Equal(t, 1, m.Get("key1"))
	assert.Equal(t, 1, m.Len())

	// Test Read non-existing key
	_, ok = m.Read("key2")
	assert.False(t, ok, "Expected to not find key 'key2'")

	// Test Delete
	m.Delete("key1")
	_, ok = m.Read("key1")
	assert.False(t, ok, "Expected to not find key 'key1' after deletion")
	assert.Equal(t, 0, m.Get("key1"))
	assert.Equal(t, 0, m.Len())

	// Test For
	m.Write("key1", 1)
	m.Write("key2", 2)
	m.For(func(k string, v int) {
		if k == "key1" {
			assert.Equal(t, 1, v, "Expected to find key 'key1' with value 1")
		} else if k == "key2" {
			assert.Equal(t, 2, v, "Expected to find key 'key2' with value 2")
		} else {
			t.Errorf("Unexpected key %s with value %d", k, v)
		}
	})
}
