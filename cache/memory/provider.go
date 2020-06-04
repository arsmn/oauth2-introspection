package memory

import (
	"context"
)

func New(cfg Config) (*Provider, error) {
	p := &Provider{
		config: cfg,
		db:     new(dict),
	}

	return p, nil
}

func (p *Provider) Get(ctx context.Context, key string) (value interface{}, ok bool) {
	val := p.db.Get(key)
	if val == nil {
		return nil, false
	}
	return val, true
}

func (p *Provider) Set(ctx context.Context, key string, value interface{}) {
	p.db.Set(key, value)
}
