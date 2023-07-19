package lrusvc

import (
	lru "github.com/hashicorp/golang-lru"
)

func NewLruCache() *lru.Cache {
	cache, _ := lru.New(100)
	return cache
}
