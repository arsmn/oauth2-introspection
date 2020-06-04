package memory

import (
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
