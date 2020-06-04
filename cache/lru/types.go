package lru

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type Config struct {
	Size       int
	Expiration time.Duration
}

type Provider struct {
	config Config
	db     *lru.Cache
}

type item struct {
	data       []byte
	expiration time.Duration
	storeTime  time.Time
}

func (i *item) isExpired() bool {
	return time.Now().UTC().After(i.storeTime.UTC().Add(i.expiration))
}
