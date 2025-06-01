package google

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/tmc/langchaingo/llms/googleai"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"
)

func (g *Google) preSendingHook(ctx context.Context, model *genai.GenerativeModel, meta googleai.PreSendingHookMetadata) {
	key := getCacheKey(model)
	if key == "" {
		return
	}

	if g.cacheNames.Get(key) == "" {
		err := g.createNewToolsCache(ctx, key, meta.Options.Model, model.Tools)
		if err != nil {
			slog.Warn("Error creating cache", "error", err)
			return
		}
	}

	model.CachedContentName = g.cacheNames.Get(key)
	model.Tools = nil

	return
}

func getCacheKey(model *genai.GenerativeModel) string {
	var funcNames []string

	for _, tool := range model.Tools {
		for _, fd := range tool.FunctionDeclarations {
			funcNames = append(funcNames, fd.Name)
		}
	}
	slices.Sort(funcNames)

	return strings.Join(funcNames, "")
}

func (g *Google) createNewToolsCache(ctx context.Context, cKey, model string, tools []*genai.Tool) error {
	cache, err := genaiCreateCachedContent(g, ctx, &genai.CachedContent{
		Expiration: genai.ExpireTimeOrTTL{
			TTL: g.cacheTTL,
		},
		Model: model,
		Tools: tools,
	})
	if err != nil {
		return err
	}

	g.cacheNames.Write(cKey, cache.Name)
	go g.startCacheRefresher(cKey, cache)

	return nil
}

func (g *Google) startCacheRefresher(cKey string, cache *genai.CachedContent) {
	go func() {
		select {
		case <-g.clientCtx.Done():
			slog.Debug("cache refresher stopped", "cacheName", cache.Name)
			return
		case <-time.After(g.cacheRefresh):
		}

		newCache, err := genaiUpdateCachedContent(g, g.clientCtx, cache, &genai.CachedContentToUpdate{
			Expiration: &genai.ExpireTimeOrTTL{
				TTL: g.cacheTTL,
			},
		})
		if err != nil {
			slog.Warn("Error refreshing cache", "cacheName", cache.Name, "error", err)
			return
		}
		g.cacheNames.Write(cKey, newCache.Name)

		go g.startCacheRefresher(cKey, newCache)
	}()
}

func (g *Google) removeAllCaches(ctx context.Context) {
	wg := &sync.WaitGroup{}
	g.cacheNames.For(func(_, name string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.removeCache(ctx, name)
		}()
	})

	wg.Wait()
}

func (g *Google) removeCache(ctx context.Context, name string) {
	if err := genaiDeleteCachedContent(g, ctx, name); err != nil {
		slog.Error("Error deleting cache", "cacheName", name, "error", err)
	} else {
		slog.Debug("Cache deleted", "cacheName", name)
	}
}

var genaiCreateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent) (*genai.CachedContent, error) {
	return g.client.GetGenaiClient().CreateCachedContent(ctx, cc)
}

var genaiUpdateCachedContent = func(g *Google, ctx context.Context, cc *genai.CachedContent, update *genai.CachedContentToUpdate) (*genai.CachedContent, error) {
	return g.client.GetGenaiClient().UpdateCachedContent(ctx, cc, update)
}

var genaiDeleteCachedContent = func(g *Google, ctx context.Context, name string) error {
	return g.client.GetGenaiClient().DeleteCachedContent(ctx, name)
}
