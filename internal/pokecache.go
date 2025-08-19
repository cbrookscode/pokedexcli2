package internal

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	interval time.Duration
	Entries  map[string]cacheEntry
	mu       sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newEntry := cacheEntry{createdAt: time.Now(), val: val}
	c.Entries[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.Entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(t *time.Ticker) {
	for {
		select {
		case <-t.C:
			c.mu.Lock()
			for key, chache_entry := range c.Entries {
				t := chache_entry.createdAt.Add(c.interval).Compare(time.Now())
				// cache duration for entry has expired
				if t != 1 {
					delete(c.Entries, key)
				}
			}
			c.mu.Unlock()
		default:
			// do nothing
		}
	}
}

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{interval: duration}
	cache.Entries = make(map[string]cacheEntry)
	timer := time.NewTicker(duration)
	go cache.reapLoop(timer)
	return cache
}

func PrintLocationsFromCache(bytes []byte) (ListofLocations, error) {
	locations := ListofLocations{}
	err := json.Unmarshal(bytes, &locations)
	if err != nil {
		return locations, fmt.Errorf("error unmarshalling data from cache into list of locations struct: %v", err)
	}

	fmt.Println("-------------------------------")
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	fmt.Println("-------------------------------")

	return locations, nil
}
