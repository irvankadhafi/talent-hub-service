package config

import "time"

const (
	EnvProduction = "production"

	DefaultDatabaseMaxIdleConns    = 3
	DefaultDatabaseMaxOpenConns    = 5
	DefaultDatabaseConnMaxLifetime = 1 * time.Hour
	DefaultDatabasePingInterval    = 1 * time.Second
	DefaultDatabaseRetryAttempts   = 3

	DefaultSessionTokenLength     = 80
	DefaultAccessTokenDuration    = 24 * time.Hour
	DefaultRefreshTokenDuration   = 24 * time.Hour * 365 // 1 year
	DefaultMaxActiveSession       = 20
	DefaultSessionDeleteBatchSize = 25
)
