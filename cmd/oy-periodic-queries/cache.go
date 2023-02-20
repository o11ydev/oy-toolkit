// Copyright The o11y toolkit Authors
// spdx-license-identifier: apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"sync"
	"time"

	"github.com/prometheus/common/model"
)

// Cache represents a simple cache with key-value pairs.
type Cache struct {
	mu      sync.Mutex
	items   map[string]cacheItem
	timeout time.Duration
}

type cacheItem struct {
	value  model.Value
	expiry time.Time
}

// NewCache creates a new cache with the specified timeout duration.
func NewCache(timeout time.Duration) *Cache {
	c := &Cache{
		items:   make(map[string]cacheItem),
		timeout: timeout,
	}
	go c.startCleanup()
	return c
}

// Get retrieves the value associated with the specified key.
func (c *Cache) Get(key string) (model.Value, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiry) {
		return nil, false
	}
	return item.value, true
}

// Set adds or updates the value associated with the specified key.
func (c *Cache) Set(key string, value model.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := cacheItem{
		value:  value,
		expiry: time.Now().Add(c.timeout),
	}
	c.items[key] = item
}

// startCleanup periodically removes expired items from the cache.
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.timeout)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mu.Lock()
		maxDelete := len(c.items) / 20
		deleted := 0
		for key, item := range c.items {
			if time.Now().After(item.expiry) {
				delete(c.items, key)
				deleted++
				if deleted > maxDelete {
					break
				}
			}
		}
		c.mu.Unlock()
	}
}
