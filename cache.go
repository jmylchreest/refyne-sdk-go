package refyne

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

// ParseCacheControl parses a Cache-Control header into directives.
func ParseCacheControl(header string) CacheControlDirectives {
	d := CacheControlDirectives{}
	if header == "" {
		return d
	}

	parts := strings.Split(strings.ToLower(header), ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case part == "no-store":
			d.NoStore = true
		case part == "no-cache":
			d.NoCache = true
		case part == "private":
			d.Private = true
		case strings.HasPrefix(part, "max-age="):
			if v, err := strconv.Atoi(part[8:]); err == nil {
				d.MaxAge = &v
			}
		case strings.HasPrefix(part, "stale-while-revalidate="):
			if v, err := strconv.Atoi(part[23:]); err == nil {
				d.StaleWhileRevalidate = &v
			}
		}
	}

	return d
}

// CreateCacheEntry creates a cache entry from a response.
// Returns nil if the response should not be cached.
func CreateCacheEntry(value any, cacheControlHeader string) *CacheEntry {
	cc := ParseCacheControl(cacheControlHeader)

	// Don't cache if no-store
	if cc.NoStore {
		return nil
	}

	// Need max-age to cache
	if cc.MaxAge == nil {
		return nil
	}

	expiresAt := time.Now().Unix() + int64(*cc.MaxAge)

	return &CacheEntry{
		Value:        value,
		ExpiresAt:    expiresAt,
		CacheControl: cc,
	}
}

// MemoryCache is an in-memory cache implementation.
type MemoryCache struct {
	store      map[string]*CacheEntry
	order      []string
	maxEntries int
	mu         sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache(maxEntries int) *MemoryCache {
	return &MemoryCache{
		store:      make(map[string]*CacheEntry),
		order:      make([]string, 0, maxEntries),
		maxEntries: maxEntries,
	}
}

// Get retrieves a cached entry by key.
func (c *MemoryCache) Get(key string) (*CacheEntry, bool) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	now := time.Now().Unix()

	// Check if expired
	if entry.ExpiresAt < now {
		// Check stale-while-revalidate
		if entry.CacheControl.StaleWhileRevalidate != nil {
			staleDeadline := entry.ExpiresAt + int64(*entry.CacheControl.StaleWhileRevalidate)
			if now < staleDeadline {
				return entry, true
			}
		}

		// Fully expired
		c.Delete(key)
		return nil, false
	}

	return entry, true
}

// Set stores an entry in the cache.
func (c *MemoryCache) Set(key string, entry *CacheEntry) {
	if entry.CacheControl.NoStore {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest if at capacity
	for len(c.store) >= c.maxEntries && len(c.order) > 0 {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.store, oldest)
	}

	// Check if key already exists
	if _, exists := c.store[key]; !exists {
		c.order = append(c.order, key)
	}

	c.store[key] = entry
}

// Delete removes an entry from the cache.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)

	// Remove from order
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

// Clear removes all entries from the cache.
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]*CacheEntry)
	c.order = make([]string, 0, c.maxEntries)
}

// Size returns the number of entries in the cache.
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}
