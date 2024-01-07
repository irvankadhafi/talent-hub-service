package db

import (
	"github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

// RedisConnectionPoolOptions options for the redis connection
type RedisConnectionPoolOptions struct {
	// Dial timeout for establishing new connections.
	// Default is 5 seconds. Only for go-redis.
	DialTimeout time.Duration

	// Enables read-only commands on slave nodes.
	ReadOnly bool

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds. Only for go-redis.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout. Only for go-redis.
	WriteTimeout time.Duration

	// Number of idle connections in the pool.
	IdleCount int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	PoolSize int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
}

var defaultRedisConnectionPoolOptions = &RedisConnectionPoolOptions{
	IdleCount:       20,
	PoolSize:        100,
	IdleTimeout:     60 * time.Second,
	MaxConnLifetime: 0,
	DialTimeout:     5 * time.Second,
	WriteTimeout:    2 * time.Second,
	ReadTimeout:     2 * time.Second,
}

// InitializeRedigoRedisConnectionPool uses redigo library to establish the redis connection pool
func InitializeRedigoRedisConnectionPool(url string, opt *RedisConnectionPoolOptions) (*redigo.Pool, error) {
	if !isValidRedisStandaloneURL(url) {
		log.Fatal("invalid redis url :", url)
	}

	options := applyRedisConnectionPoolOptions(opt)

	return &redigo.Pool{
		MaxIdle:     options.IdleCount,
		MaxActive:   options.PoolSize,
		IdleTimeout: options.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		MaxConnLifetime: options.MaxConnLifetime,
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Wait: true, // wait for connection available when maxActive is reached
	}, nil
}

func isValidRedisStandaloneURL(url string) bool {
	_, err := redis.ParseURL(url)
	if err != nil {
		log.Error(err)
	}

	return err == nil
}

func applyRedisConnectionPoolOptions(opt *RedisConnectionPoolOptions) *RedisConnectionPoolOptions {
	if opt != nil {
		return opt
	}

	return defaultRedisConnectionPoolOptions
}
