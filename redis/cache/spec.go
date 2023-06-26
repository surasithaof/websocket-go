package cache

import (
	"context"
	"time"
)

//go:generate mockgen -source=./spec.go -destination=./mocks/cache.go -package=mocks "surasithaof/websocket-go" Cache

type Cache interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
