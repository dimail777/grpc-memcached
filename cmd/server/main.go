package main

import (
	"internal/app"
	"internal/memcache"
	"os"
	"strconv"
	"time"
)

func main() {
	mode := os.Getenv("MODE")
	if mode == "memcached" {
		app.RunMemcached(initMemcachedEnv())
	} else {
		app.RunGoCached(initGoCachedEnv())
	}
}

const (
	defaultConcurrency             = 16
	defaultHost                    = "localhost"
	defaultPort              int32 = 11211
	defaultTcpTimeout              = 5 * time.Second
	defaultMaxConnections          = 10
	defaultMaxWaitConnection       = 30 * time.Second
)

func initGoCachedEnv() int {
	return getEnvOrDefaultInt("GO_CACHED_CONCURRENCY", defaultConcurrency)
}

func initMemcachedEnv() (*memcache.ConnConfig, *memcache.PoolConfig) {
	connConfig := memcache.ConnConfig{
		Host:       getEnvOrDefault("MEMCACHED_HOST", defaultHost),
		TcpTimeout: getEnvOrDefaultDuration("MEMCACHED_TCP_TIMEOUT", defaultTcpTimeout),
	}
	connConfig.Port = int32(getEnvOrDefaultInt("MEMCACHED_PORT", int(defaultPort)))

	poolConfig := memcache.PoolConfig{
		MaxConnections:    getEnvOrDefaultInt("MEMCACHED_MAX_CONNECTIONS", defaultMaxConnections),
		MaxWaitConnection: getEnvOrDefaultDuration("MEMCACHED_MAX_WAIT_CONNECTION", defaultMaxWaitConnection),
	}
	return &connConfig, &poolConfig
}
func getEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}
