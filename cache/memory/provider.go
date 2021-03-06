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
	item.expiration = 0
	item.storeTime = time.Time{}

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
	if item.isExpired() {
		p.db.Del(key)
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

	p.db.Set(key, item)

	return nil
}
