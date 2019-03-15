package cache

import (
	"time"
)

type CacheItem struct {
	value interface{} // cached value
	btl   *uint64     // expiration block id
	ttl   *time.Time  // expiration time
}

func (c *CacheItem) IsExpired(currentBlock uint64) bool {
	if c.btl != nil && currentBlock <= *c.btl {
		return false
	}

	if c.ttl != nil && time.Now().Before(*c.ttl) {
		return false
	}

	return true
}
