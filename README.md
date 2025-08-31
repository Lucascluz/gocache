# GoCache

A simple, fast in-memory cache for Go with HTTP server capabilities.

## Features

- Thread-safe in-memory cache with TTL support
- Optional HTTP server with GET/SET/DELETE endpoints  
- Configurable cleanup intervals
- Zero external dependencies (uses only Go standard library)

## Installation

```bash
go get github.com/Lucascluz/gocache
```

## Quick Start

### Basic Cache Usage

```go
package main

import (
    "time"
    "github.com/Lucascluz/gocache/pkg/cache"
)

func main() {
    // Create cache with 5-minute cleanup interval
    store := cache.New(&cache.Config{
        CleanupInterval: 5 * time.Minute,
    })
    
    // Set values
    store.Set("key1", "value1")
    store.SetWithTTL("key2", "value2", 30*time.Second)
    
    // Get values
    if val, ok := store.Get("key1"); ok {
        println("Found:", val.(string))
    }
}
```

### HTTP Server

```go
package main

import (
    "context"
    "time"
    
    "github.com/Lucascluz/gocache/pkg/cache"
    "github.com/Lucascluz/gocache/pkg/gocache"
)

func main() {
    cfg := gocache.Config{
        CacheConfig: cache.Config{
            CleanupInterval: 5 * time.Minute,
        },
        HttpConfig: http.Config{
            Enabled: true,
            Port:    8080,
        },
    }
    
    server := gocache.New(&cfg)
    
    // Access cache directly
    server.Cache().Set("test", "value")
    
    // Start HTTP server (blocks)
    server.ListenAndServe()
}
```

### HTTP Endpoints

- `GET /get` - Get value by key (header: `key` or query: `?key=mykey`)
- `POST /set` - Set value (header: `key`, body: value, optional header: `ttl-seconds`)
- `DELETE /delete` - Delete key (header: `key` or query: `?key=mykey`)
- `GET /health` - Health check

## API Reference

See the [pkg.go.dev documentation](https://pkg.go.dev/github.com/Lucascluz/gocache) for detailed API reference.

## License

MIT - see [LICENSE](LICENSE) file.
