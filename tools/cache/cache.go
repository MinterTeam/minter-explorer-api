package cache

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"sync"
	"time"
)

type ExplorerCache struct {
	lastBlock blocks.Resource
	items     *sync.Map
}

// cache constructor
func NewCache(lastBlock models.Block) *ExplorerCache {
	cache := &ExplorerCache{
		lastBlock: new(blocks.Resource).Transform(lastBlock).(blocks.Resource),
		items:     new(sync.Map),
	}

	return cache
}

// create new cache item
func (c *ExplorerCache) newCacheItem(value interface{}, ttl interface{}) *CacheItem {
	switch t := ttl.(type) {
	case time.Duration:
		ttl := time.Now().Add(t * time.Second)
		return &CacheItem{value: value, ttl: &ttl}
	case int:
		ttl := c.lastBlock.ID + uint64(t)
		return &CacheItem{value: value, btl: &ttl}
	}

	panic("Invalid cache ttl type.")
}

// get or store value from cache
func (c *ExplorerCache) Get(key interface{}, callback func() interface{}, ttl interface{}) interface{} {
	v, ok := c.items.Load(key)
	if ok {
		item := v.(*CacheItem)
		if !item.IsExpired(c.lastBlock.ID) {
			return item.value
		}
	}

	return c.Store(key, callback(), ttl)
}

// get or store value from cache
func (c *ExplorerCache) ExecuteOrGet(key interface{}, callback func() interface{}, ttl interface{}, executeRequired bool) interface{} {
	if executeRequired {
		return callback()
	}

	v, ok := c.items.Load(key)
	if ok {
		item := v.(*CacheItem)
		if !item.IsExpired(c.lastBlock.ID) {
			return item.value
		}
	}

	return c.Store(key, callback(), ttl)
}

// save value to cache
func (c *ExplorerCache) Store(key interface{}, value interface{}, ttl interface{}) interface{} {
	c.items.Store(key, c.newCacheItem(value, ttl))
	return value
}

// loop for checking items expiration
func (c *ExplorerCache) ExpirationCheck() {
	c.items.Range(func(key, value interface{}) bool {
		item := value.(*CacheItem)
		if item.IsExpired(c.lastBlock.ID) {
			c.items.Delete(key)
		}

		return true
	})
}

// set new last block id
func (c *ExplorerCache) SetLastBlock(block blocks.Resource) {
	c.lastBlock = block
	// clean expired items
	go c.ExpirationCheck()
}

// Get latest explorer block
func (c *ExplorerCache) GetLastBlock() blocks.Resource {
	return c.lastBlock
}

// update last block id by ws data
func (c *ExplorerCache) OnNewBlock(block blocks.Resource) {
	c.SetLastBlock(block)
}
