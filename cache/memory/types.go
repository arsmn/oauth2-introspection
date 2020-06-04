package memory

import (
	"time"

	"github.com/savsgio/dictpool"
)

// Config provider settings
type Config struct{}

// Provider backend manager
type Provider struct {
	config Config
	db     *dict
}

type dict struct {
	dictpool.Dict
}

type item struct {
	data           []byte
	lastActiveTime int64
	expiration     time.Duration
}
