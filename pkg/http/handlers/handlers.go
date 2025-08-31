package handlers

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	stdhttp "net/http"

	"github.com/Lucascluz/gocache/pkg/cache"
)

// GET returns an http.HandlerFunc that reads the "key" from header or query and fetches from cache.
func GET(store *cache.Cache) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		key := r.Header.Get("key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}
		if key == "" {
			stdhttp.Error(w, "missing key", stdhttp.StatusBadRequest)
			return
		}

		val, ok := store.Get(key)
		if !ok {
			stdhttp.Error(w, "not found", stdhttp.StatusNotFound)
			return
		}

		// Best-effort string representation
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(strings.TrimSpace(toString(val))))
	}
}

// SET returns an http.HandlerFunc that writes a value to the cache.
// Key is taken from header or query. Body is used as the value (string).
// Optional header: ttl-seconds to set a TTL.
func SET(store *cache.Cache) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		key := r.Header.Get("key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}
		if key == "" {
			stdhttp.Error(w, "missing key", stdhttp.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			stdhttp.Error(w, err.Error(), stdhttp.StatusBadRequest)
			return
		}
		_ = r.Body.Close()

		ttlHeader := r.Header.Get("ttl-seconds")
		if ttlHeader != "" {
			if secs, err := strconv.ParseInt(ttlHeader, 10, 64); err == nil && secs > 0 {
				_ = store.SetWithTTL(key, string(body), time.Duration(secs)*time.Second)
				w.WriteHeader(stdhttp.StatusNoContent)
				return
			}
		}

		_ = store.Set(key, string(body))
		w.WriteHeader(stdhttp.StatusNoContent)
	}
}

// DELETE returns an http.HandlerFunc that deletes a key from the cache.
func DELETE(store *cache.Cache) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		key := r.Header.Get("key")
		if key == "" {
			key = r.URL.Query().Get("key")
		}
		if key == "" {
			stdhttp.Error(w, "missing key", stdhttp.StatusBadRequest)
			return
		}
		if store.Delete(key) {
			w.WriteHeader(stdhttp.StatusNoContent)
			return
		}
		stdhttp.Error(w, "not found", stdhttp.StatusNotFound)
	}
}

// toString provides a simple string conversion for common types.
func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", t)
	}
}
