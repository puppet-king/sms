// Copyright 2022 The wangkai. ALL rights reserved.

/*
Package models
*/
package models

import (
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

var (
	once  sync.Once
	cache *Cache
)

func NewCache() *Cache {
	once.Do(func() {
		cache = &Cache{
			data: make(map[string]interface{}),
		}
	})

	return cache
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

// Set exp 存在 bug 无法修改
func (c *Cache) Set(key string, val interface{}, exp time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
	if exp > 0 {
		time.AfterFunc(exp, func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			delete(c.data, key)
		})
	}
}

// SetAllowList 设置允许列表
func (c *Cache) SetAllowList() {
	newCache := NewCache()
	var allowOpenidList = make(map[string]bool)

	info := TableUser{}
	user := info.List()
	for _, v := range user {
		allowOpenidList[v.UserId] = true
	}

	newCache.Set("sms:user:allowList", allowOpenidList, 0)
}

// GetAllowList 获取允许列表
func (c *Cache) GetAllowList() map[string]bool {
	newCache := NewCache()
	if data, ok := newCache.Get("sms:user:allowList"); ok {
		if allowOpenidList, ok := data.(map[string]bool); ok {
			return allowOpenidList
		}
	}

	return map[string]bool{}
}

// GetLastLoginInfo 获取允许列表
func (c *Cache) GetLastLoginInfo() map[string]string {
	newCache := NewCache()
	if data, ok := newCache.Get("sms:user:allowList"); ok {
		if allowOpenidList, ok := data.(map[string]string); ok {
			return allowOpenidList
		}
	}

	return map[string]string{}
}
