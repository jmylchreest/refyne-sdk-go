package refyne

import (
	"testing"
	"time"
)

func TestParseCacheControl(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   CacheControlDirectives
	}{
		{
			name:   "empty header",
			header: "",
			want:   CacheControlDirectives{},
		},
		{
			name:   "no-store",
			header: "no-store",
			want:   CacheControlDirectives{NoStore: true},
		},
		{
			name:   "no-cache",
			header: "no-cache",
			want:   CacheControlDirectives{NoCache: true},
		},
		{
			name:   "private",
			header: "private",
			want:   CacheControlDirectives{Private: true},
		},
		{
			name:   "max-age",
			header: "max-age=3600",
			want:   CacheControlDirectives{MaxAge: intPtr(3600)},
		},
		{
			name:   "multiple directives",
			header: "private, max-age=3600, stale-while-revalidate=60",
			want: CacheControlDirectives{
				Private:              true,
				MaxAge:               intPtr(3600),
				StaleWhileRevalidate: intPtr(60),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCacheControl(tt.header)

			if got.NoStore != tt.want.NoStore {
				t.Errorf("NoStore = %v, want %v", got.NoStore, tt.want.NoStore)
			}
			if got.NoCache != tt.want.NoCache {
				t.Errorf("NoCache = %v, want %v", got.NoCache, tt.want.NoCache)
			}
			if got.Private != tt.want.Private {
				t.Errorf("Private = %v, want %v", got.Private, tt.want.Private)
			}
			if !intPtrEqual(got.MaxAge, tt.want.MaxAge) {
				t.Errorf("MaxAge = %v, want %v", got.MaxAge, tt.want.MaxAge)
			}
			if !intPtrEqual(got.StaleWhileRevalidate, tt.want.StaleWhileRevalidate) {
				t.Errorf("StaleWhileRevalidate = %v, want %v", got.StaleWhileRevalidate, tt.want.StaleWhileRevalidate)
			}
		})
	}
}

func TestCreateCacheEntry(t *testing.T) {
	t.Run("returns nil for no-store", func(t *testing.T) {
		entry := CreateCacheEntry("value", "no-store")
		if entry != nil {
			t.Error("expected nil for no-store")
		}
	})

	t.Run("returns nil without max-age", func(t *testing.T) {
		entry := CreateCacheEntry("value", "private")
		if entry != nil {
			t.Error("expected nil without max-age")
		}
	})

	t.Run("creates entry with max-age", func(t *testing.T) {
		now := time.Now().Unix()
		entry := CreateCacheEntry("value", "max-age=3600")

		if entry == nil {
			t.Fatal("expected non-nil entry")
		}
		if entry.ExpiresAt < now+3600 {
			t.Error("expires_at too early")
		}
	})
}

func TestMemoryCache(t *testing.T) {
	t.Run("stores and retrieves", func(t *testing.T) {
		cache := NewMemoryCache(10)
		entry := &CacheEntry{
			Value:     "test",
			ExpiresAt: time.Now().Unix() + 3600,
		}

		cache.Set("key", entry)
		got, ok := cache.Get("key")

		if !ok {
			t.Error("expected entry to be found")
		}
		if got.Value != "test" {
			t.Errorf("got %v, want test", got.Value)
		}
	})

	t.Run("returns false for missing", func(t *testing.T) {
		cache := NewMemoryCache(10)
		_, ok := cache.Get("nonexistent")
		if ok {
			t.Error("expected false for nonexistent key")
		}
	})

	t.Run("expires entries", func(t *testing.T) {
		cache := NewMemoryCache(10)
		entry := &CacheEntry{
			Value:     "test",
			ExpiresAt: time.Now().Unix() - 1,
		}

		cache.Set("key", entry)
		_, ok := cache.Get("key")

		if ok {
			t.Error("expected expired entry to not be found")
		}
	})

	t.Run("evicts oldest", func(t *testing.T) {
		cache := NewMemoryCache(2)
		future := time.Now().Unix() + 3600

		cache.Set("key1", &CacheEntry{Value: "v1", ExpiresAt: future})
		cache.Set("key2", &CacheEntry{Value: "v2", ExpiresAt: future})
		cache.Set("key3", &CacheEntry{Value: "v3", ExpiresAt: future})

		_, ok := cache.Get("key1")
		if ok {
			t.Error("expected key1 to be evicted")
		}

		_, ok = cache.Get("key2")
		if !ok {
			t.Error("expected key2 to exist")
		}

		_, ok = cache.Get("key3")
		if !ok {
			t.Error("expected key3 to exist")
		}
	})

	t.Run("delete", func(t *testing.T) {
		cache := NewMemoryCache(10)
		cache.Set("key", &CacheEntry{Value: "test", ExpiresAt: time.Now().Unix() + 3600})
		cache.Delete("key")

		_, ok := cache.Get("key")
		if ok {
			t.Error("expected key to be deleted")
		}
	})

	t.Run("clear", func(t *testing.T) {
		cache := NewMemoryCache(10)
		cache.Set("key1", &CacheEntry{Value: "v1", ExpiresAt: time.Now().Unix() + 3600})
		cache.Set("key2", &CacheEntry{Value: "v2", ExpiresAt: time.Now().Unix() + 3600})
		cache.Clear()

		if cache.Size() != 0 {
			t.Errorf("expected size 0, got %d", cache.Size())
		}
	})
}

func intPtr(i int) *int {
	return &i
}

func intPtrEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
