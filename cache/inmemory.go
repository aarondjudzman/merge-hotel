package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func NewInMemoryCache() *cache.Cache {
	c := cache.New(1*time.Minute, 2*time.Minute)

	return c
}
