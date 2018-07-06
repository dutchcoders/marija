package server

import (
	_ "log"
	"sync"

	"github.com/dutchcoders/marija/server/datasources"
)

type ItemCache struct {
	sync.Map
}

func (ic *ItemCache) Store(key string, items []datasources.Item) {
	ic.Map.Store(key, items)
}

func (is *ItemCache) LoadOrStore(key string, value []datasources.Item) ([]datasources.Item, bool) {
	actual, ok := is.Map.LoadOrStore(key, value)
	if !ok {
		return nil, false
	}

	return actual.([]datasources.Item), true
}

func (is *ItemCache) Load(key string) ([]datasources.Item, bool) {
	value, ok := is.Map.Load(key)
	if !ok {
		return nil, false
	}

	return value.([]datasources.Item), true
}
