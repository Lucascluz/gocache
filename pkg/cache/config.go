package cache

import "time"

type Config struct {
	CleanupInterval time.Duration
	MaxSize         int64
}
