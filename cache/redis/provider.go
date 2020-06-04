package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var errConfigAddrEmpty = errors.New("Config Addr must not be empty")

// New returns a new configured redis provider
func New(cfg Config) (*Provider, error) {
	if cfg.Addr == "" {
		return nil, errConfigAddrEmpty
	}

	db := redis.NewClient(&redis.Options{
		Network:            cfg.Network,
		Addr:               cfg.Addr,
		Password:           cfg.Password,
		DB:                 cfg.DB,
		MaxRetries:         cfg.MaxRetries,
		MinRetryBackoff:    cfg.MinRetryBackoff,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		MinIdleConns:       cfg.MinIdleConns,
		MaxConnAge:         cfg.MaxConnAge,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		TLSConfig:          cfg.TLSConfig,
		Limiter:            cfg.Limiter,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("Redis connection error: %v", err)
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

func (p *Provider) Get(key string) ([]byte, error) {
	val, err := p.db.Get(context.Background(), key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return val, nil
}

func (p *Provider) Set(key string, data []byte, expiration time.Duration) error {
	return p.db.Set(context.Background(), key, data, expiration).Err()
}
