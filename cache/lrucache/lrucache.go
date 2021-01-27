package lrucache

import (
	"sync"

	"github.com/golang/groupcache/lru"
)

type LRUCache struct {
	data *lru.Cache
	m    *sync.Mutex
}
