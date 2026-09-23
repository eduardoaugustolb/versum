package config

import "errors"

var (
	ErrRedisURLNotSet = errors.New("redis url not set")
)

const (
	DefaultRedisURLKey = "REDIS_URL"
)
