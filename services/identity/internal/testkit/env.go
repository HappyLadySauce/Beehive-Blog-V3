package testkit

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// testPostgresEnv describes fallback PostgreSQL settings for integration tests.
// testPostgresEnv 描述集成测试的 PostgreSQL fallback 配置。
type testPostgresEnv struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// testRedisEnv describes fallback Redis settings for integration tests.
// testRedisEnv 描述集成测试的 Redis fallback 配置。
type testRedisEnv struct {
	Host     string
	Port     int
	Password string
	Database int
}

// loadPostgresEnv loads fallback PostgreSQL settings from environment variables.
// loadPostgresEnv 从环境变量加载 PostgreSQL fallback 配置。
func loadPostgresEnv() (*testPostgresEnv, bool, error) {
	host := strings.TrimSpace(os.Getenv("BEEHIVE_TEST_PG_HOST"))
	if host == "" {
		return nil, false, nil
	}

	port, err := readIntEnv("BEEHIVE_TEST_PG_PORT", 5432)
	if err != nil {
		return nil, false, err
	}

	cfg := &testPostgresEnv{
		Host:     host,
		Port:     port,
		User:     strings.TrimSpace(defaultString(os.Getenv("BEEHIVE_TEST_PG_USER"), "postgres")),
		Password: os.Getenv("BEEHIVE_TEST_PG_PASSWORD"),
		DBName:   strings.TrimSpace(defaultString(os.Getenv("BEEHIVE_TEST_PG_DBNAME"), "postgres")),
		SSLMode:  strings.TrimSpace(defaultString(os.Getenv("BEEHIVE_TEST_PG_SSLMODE"), "disable")),
		TimeZone: strings.TrimSpace(defaultString(os.Getenv("BEEHIVE_TEST_PG_TIMEZONE"), "Asia/Shanghai")),
	}

	return cfg, true, nil
}

// loadRedisEnv loads fallback Redis settings from environment variables.
// loadRedisEnv 从环境变量加载 Redis fallback 配置。
func loadRedisEnv() (*testRedisEnv, bool, error) {
	host := strings.TrimSpace(os.Getenv("BEEHIVE_TEST_REDIS_HOST"))
	if host == "" {
		return nil, false, nil
	}

	port, err := readIntEnv("BEEHIVE_TEST_REDIS_PORT", 6379)
	if err != nil {
		return nil, false, err
	}
	database, err := readIntEnv("BEEHIVE_TEST_REDIS_DB", 0)
	if err != nil {
		return nil, false, err
	}

	cfg := &testRedisEnv{
		Host:     host,
		Port:     port,
		Password: os.Getenv("BEEHIVE_TEST_REDIS_PASSWORD"),
		Database: database,
	}

	return cfg, true, nil
}

// readIntEnv reads an integer environment variable with a default value.
// readIntEnv 读取整数环境变量，并在缺失时返回默认值。
func readIntEnv(key string, defaultValue int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}

	return value, nil
}

// defaultString returns a fallback value for empty strings.
// defaultString 在字符串为空时返回默认值。
func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
