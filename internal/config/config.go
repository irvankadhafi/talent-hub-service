package config

import (
	"fmt"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warningf("%v", err)
	}
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return viper.GetString("ports.http")
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		DatabaseUsername(),
		DatabasePassword(),
		DatabaseHost(),
		DatabaseName(),
		DatabaseSSLMode())
}

// DatabaseHost :nodoc:
func DatabaseHost() string {
	return viper.GetString("postgres.host")
}

// DatabaseName :nodoc:
func DatabaseName() string {
	return viper.GetString("postgres.database")
}

// DatabaseUsername :nodoc:
func DatabaseUsername() string {
	return viper.GetString("postgres.username")
}

// DatabasePassword :nodoc:
func DatabasePassword() string {
	return viper.GetString("postgres.password")
}

// DatabaseSSLMode :nodoc:
func DatabaseSSLMode() string {
	if viper.IsSet("postgres.sslmode") {
		return viper.GetString("postgres.sslmode")
	}
	return "disable"
}

// DatabasePingInterval :nodoc:
func DatabasePingInterval() time.Duration {
	if viper.GetInt("postgres.ping_interval") <= 0 {
		return DefaultDatabasePingInterval
	}
	return time.Duration(viper.GetInt("postgres.ping_interval")) * time.Millisecond
}

// DatabaseRetryAttempts :nodoc:
func DatabaseRetryAttempts() float64 {
	if viper.GetInt("postgres.retry_attempts") > 0 {
		return float64(viper.GetInt("postgres.retry_attempts"))
	}
	return DefaultDatabaseRetryAttempts
}

// DatabaseMaxIdleConns :nodoc:
func DatabaseMaxIdleConns() int {
	if viper.GetInt("postgres.max_idle_conns") <= 0 {
		return DefaultDatabaseMaxIdleConns
	}
	return viper.GetInt("postgres.max_idle_conns")
}

// DatabaseMaxOpenConns :nodoc:
func DatabaseMaxOpenConns() int {
	if viper.GetInt("postgres.max_open_conns") <= 0 {
		return DefaultDatabaseMaxOpenConns
	}
	return viper.GetInt("postgres.max_open_conns")
}

// DatabaseConnMaxLifetime :nodoc:
func DatabaseConnMaxLifetime() time.Duration {
	if !viper.IsSet("postgres.conn_max_lifetime") {
		return DefaultDatabaseConnMaxLifetime
	}
	return time.Duration(viper.GetInt("postgres.conn_max_lifetime")) * time.Millisecond
}

func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

func RedisWorkerHost() string {
	return viper.GetString("redis.worker_host")
}

func RedisDialTimeout() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("redis.dial_timeout"), 5*time.Second)
}

func RedisWriteTimeout() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("redis.write_timeout"), 2*time.Second)
}

func RedisReadTimeout() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("redis.read_timeout"), 2*time.Second)
}

func RedisMaxIdleConn() int {
	return utils.ValueOrDefault[int](utils.StringToInt[int](viper.GetString("redis.max_idle_conn")), 20)
}

func RedisMaxActiveConn() int {
	return utils.ValueOrDefault[int](utils.StringToInt[int](viper.GetString("redis.max_active_conn")), 50)
}

// DisableCaching :nodoc:
func DisableCaching() bool {
	return viper.GetBool("disable_caching")
}

// SessionDeleteBatchSize get max deletion limit on each old sess deletion loop
func SessionDeleteBatchSize() int {
	cfg := viper.GetInt("session.deletion_batch_size")

	if cfg <= 0 {
		return DefaultSessionDeleteBatchSize
	}

	return cfg
}

// AccessTokenDuration get access token increment duration in hour
func AccessTokenDuration() time.Duration {
	cfg := viper.GetString("session.access_token_duration")
	return utils.ParseDurationWithDefault(cfg, DefaultAccessTokenDuration)
}

// RefreshTokenDuration get refresh token increment duration in hour
func RefreshTokenDuration() time.Duration {
	cfg := viper.GetString("session.refresh_token_duration")
	return utils.ParseDurationWithDefault(cfg, DefaultRefreshTokenDuration)
}
