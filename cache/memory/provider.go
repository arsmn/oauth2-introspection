package memory

import (
	"sync"
	"time"
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
	item.lastActiveTime = 0
	item.expiration = 0

	itemPool.Put(item)
}

func New(cfg Config) (*Provider, error) {
	p := &Provider{
		config: cfg,
		db:     new(dict),
	}

	return p, nil
}

func (p *Provider) Get(key string) ([]byte, error) {
	val := p.db.Get(key)
	if val == nil {
		return nil, nil
	}

	item := val.(*item)
	return item.data, nil
}

func (p *Provider) Set(key string, data []byte, expiration time.Duration) error {
	item := acquireItem()
	item.data = data
	item.lastActiveTime = time.Now().UnixNano()
	item.expiration = expiration

	p.db.Set(key, item)

	return nil
}
