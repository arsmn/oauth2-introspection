package memory

import (
	"time"

	"github.com/savsgio/dictpool"
)

// Config provider settings
type Config struct {
	Expiration time.Duration
}

// Provider backend manager
type Provider struct {
	config Config
	db     *dict
}

type dict struct {
	dictpool.Dict
}

type item struct {
	data       []byte
	expiration time.Duration
	storeTime  time.Time
}

func (i *item) isExpired() bool {
	return time.Now().UTC().After(i.storeTime.UTC().Add(i.expiration))
}
