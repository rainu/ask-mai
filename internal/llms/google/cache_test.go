package google

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"sync/atomic"
	"testing"
	"time"
)

func TestCaching_NoCacheNeeded(t *testing.T) {
	toTest := &Google{}
	testModel := &genai.GenerativeModel{
		Tools: make([]*genai.Tool, 0),
	}
	testMeta := googleai.PreSendingHookMetadata{}

	toTest.preSendingHook(t.Context(), testModel, testMeta)

	assert.Equal(t, "", testModel.CachedContentName)
	assert.NotNil(t, testModel.Tools)
	assert.Equal(t, 0, toTest.cacheNames.Len())
}

func TestCaching_ErrorWhileCreating(t *testing.T) {
	toTest := &Google{
		cacheTTL: 5 * time.Minute,
	}

	var err error
	toTest.client, err = googleai.New(t.Context(), googleai.WithAPIKey("API_KEY"))
	require.NoError(t, err)

	testModel := &genai.GenerativeModel{}
	for i := 0; i < 10; i++ {
		testModel.Tools = append(testModel.Tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{Name: fmt.Sprintf("func%d", i)},
			},
		})
	}

	testMeta := googleai.PreSendingHookMetadata{
		Options: llms.CallOptions{
			Model: "test-model",
		},
	}

	of := genaiCreateCachedContent
	defer func() { genaiCreateCachedContent = of }()
	genaiCreateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent) (*genai.CachedContent, error) {
		assert.Same(t, toTest, g)
		assert.Same(t, t.Context(), ctx)
		assert.Equal(t, genai.CachedContent{
			Expiration: genai.ExpireTimeOrTTL{
				TTL: toTest.cacheTTL,
			},
			Model: testMeta.Options.Model,
			Tools: testModel.Tools,
		}, *cc)
		return nil, fmt.Errorf("test error")
	}

	toTest.preSendingHook(t.Context(), testModel, testMeta)

	assert.Equal(t, "", testModel.CachedContentName)
	assert.NotNil(t, testModel.Tools)
	assert.Equal(t, 0, toTest.cacheNames.Len())
}

func TestCaching_Success(t *testing.T) {
	toTest := &Google{
		cacheTTL:  5 * time.Minute,
		clientCtx: t.Context(),
	}

	var err error
	toTest.client, err = googleai.New(t.Context(), googleai.WithAPIKey("API_KEY"))
	require.NoError(t, err)

	testModel := &genai.GenerativeModel{}
	for i := 0; i < 10; i++ {
		testModel.Tools = append(testModel.Tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{Name: fmt.Sprintf("func%d", i)},
			},
		})
	}

	testMeta := googleai.PreSendingHookMetadata{
		Options: llms.CallOptions{
			Model: "test-model",
		},
	}

	of := genaiCreateCachedContent
	defer func() { genaiCreateCachedContent = of }()
	genaiCreateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent) (*genai.CachedContent, error) {
		assert.Same(t, toTest, g)
		assert.Same(t, t.Context(), ctx)
		assert.Equal(t, genai.CachedContent{
			Expiration: genai.ExpireTimeOrTTL{
				TTL: toTest.cacheTTL,
			},
			Model: testMeta.Options.Model,
			Tools: testModel.Tools,
		}, *cc)
		return &genai.CachedContent{Name: "test-cache"}, nil
	}

	toTest.preSendingHook(t.Context(), testModel, testMeta)

	assert.Equal(t, "test-cache", testModel.CachedContentName)
	assert.Nil(t, testModel.Tools)
	assert.Equal(t, 1, toTest.cacheNames.Len())
}

func TestCaching_Refreshing(t *testing.T) {
	toTest := &Google{
		cacheRefresh: 50 * time.Millisecond,
		clientCtx:    t.Context(),
	}
	toTest.cacheTTL = toTest.cacheRefresh * 3 // ensure that the cache is refreshed at least twice

	var err error
	toTest.client, err = googleai.New(t.Context(), googleai.WithAPIKey("API_KEY"))
	require.NoError(t, err)

	testModel := &genai.GenerativeModel{
		Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{Name: "funcA"}}}},
	}

	testMeta := googleai.PreSendingHookMetadata{
		Options: llms.CallOptions{
			Model: "test-model",
		},
	}

	ogccc := genaiCreateCachedContent
	defer func() { genaiCreateCachedContent = ogccc }()
	genaiCreateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent) (*genai.CachedContent, error) {
		result := *cc
		result.Name = "test-cache"

		return &result, nil
	}

	ogucc := genaiUpdateCachedContent
	defer func() { genaiUpdateCachedContent = ogucc }()
	genaiUpdateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent, update *genai.CachedContentToUpdate) (*genai.CachedContent, error) {
		assert.Same(t, toTest, g)

		if cc.Name == "test-cache" {
			//first cache refresh
			assert.Equal(t, genai.CachedContent{
				Name: "test-cache",
				Expiration: genai.ExpireTimeOrTTL{
					TTL: toTest.cacheTTL,
				},
				Model: testMeta.Options.Model,
				Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{Name: "funcA"}}}},
			}, *cc)
		} else {
			assert.Equal(t, genai.CachedContent{
				Name: "test-cache-update",
				Expiration: genai.ExpireTimeOrTTL{
					TTL: toTest.cacheTTL,
				},
				Model: testMeta.Options.Model,
				Tools: []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{{Name: "funcA"}}}},
			}, *cc)
		}

		assert.Equal(t, genai.CachedContentToUpdate{
			Expiration: &genai.ExpireTimeOrTTL{
				TTL: toTest.cacheTTL,
			},
		}, *update)

		result := *cc
		result.Name = "test-cache-update"

		return &result, nil
	}

	toTest.preSendingHook(t.Context(), testModel, testMeta)

	//wait for refreshing to kick in
	select {
	case <-t.Context().Done():
		t.Fail()
		return
	case <-time.After(toTest.cacheTTL):
	}

	assert.Equal(t, "test-cache", testModel.CachedContentName)
	assert.Equal(t, 1, toTest.cacheNames.Len())
	toTest.cacheNames.For(func(_ string, cacheName string) {
		assert.Equal(t, "test-cache-update", cacheName)
	})
}

func TestCaching_Cleanup(t *testing.T) {
	toTest := &Google{}
	toTest.clientCtx, toTest.clientCancel = context.WithCancel(context.Background())

	var err error
	toTest.client, err = googleai.New(t.Context(), googleai.WithAPIKey("API_KEY"))
	require.NoError(t, err)

	singleEntryDuration := 5 * time.Millisecond
	called := atomic.Int32{}

	ogdcc := genaiDeleteCachedContent
	defer func() { genaiDeleteCachedContent = ogdcc }()
	genaiDeleteCachedContent = func(g *Google, ctx context.Context, name string) error {
		assert.Same(t, toTest, g)
		assert.NotSame(t, toTest.clientCtx, ctx)

		called.Add(1)
		time.Sleep(singleEntryDuration)
		return nil
	}

	entries := 100
	for i := 0; i < entries; i++ {
		toTest.cacheNames.Write(fmt.Sprintf("%d", i), "test-cache")
	}
	require.Equal(t, entries, toTest.cacheNames.Len(), "Cache should have the expected number of entries")

	// close should trigger cleanup
	startTime := time.Now()
	toTest.Close()
	duration := time.Since(startTime)

	assert.LessOrEqual(t, duration, (time.Duration(entries)*singleEntryDuration)/2)
	assert.Error(t, toTest.clientCtx.Err(), "Context should be cancelled after Close")
	assert.Equal(t, int32(entries), called.Load(), "All cache entries should have been deleted")
}

func Test_getCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		model    *genai.GenerativeModel
		expected string
	}{
		{
			name:     "No tools",
			model:    &genai.GenerativeModel{Tools: nil},
			expected: "",
		},
		{
			name: "Single tool with one function",
			model: &genai.GenerativeModel{
				Tools: []*genai.Tool{{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						{Name: "func1"},
					},
				}},
			},
			expected: "func1",
		},
		{
			name: "Multiple tools with multiple functions",
			model: &genai.GenerativeModel{
				Tools: []*genai.Tool{
					{
						FunctionDeclarations: []*genai.FunctionDeclaration{
							{Name: "funcB"},
							{Name: "funcA"},
						},
					},
					{
						FunctionDeclarations: []*genai.FunctionDeclaration{
							{Name: "funcC"},
						},
					},
				},
			},
			expected: "funcAfuncBfuncC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getCacheKey(tt.model))
		})
	}
}
