package lru

import (
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

var itemPool = &sync.Pool{
	New: func() interface{} {
		return new(item)
	},
}

func acquireItem() *item {
	return itemPool.Get().(*item)
}

func releaseItem(item *item) {
	item.data = item.data[:0]
	item.expiration = 0
	item.storeTime = time.Time{}

	itemPool.Put(item)
}

func New(cfg Config) (*Provider, error) {
	db, err := lru.New(cfg.Size)
	if err != nil {
		return nil, fmt.Errorf("LRU initialization error: %v", err)
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

func (p *Provider) Get(key string) ([]byte, error) {
	val, ok := p.db.Get(key)
	if !ok {
		return nil, nil
	}

	item := val.(*item)
	if item.isExpired() {
		p.db.Remove(key)
		releaseItem(item)
		return nil, nil
	}

	return item.data, nil
}

func (p *Provider) Set(key string, data []byte) error {
	item := acquireItem()
	item.data = data
	item.expiration = p.config.Expiration
	item.storeTime = time.Now().UTC()

	p.db.Add(key, item)
	return nil
}
