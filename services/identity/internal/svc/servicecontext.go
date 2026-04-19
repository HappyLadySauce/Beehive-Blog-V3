package svc

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	SQLDB  *sql.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("identity 配置校验失败: %w", err))
	}

	db, sqlDB, err := newPostgres(c.Postgres)
	if err != nil {
		panic(fmt.Errorf("初始化 PostgreSQL 失败: %w", err))
	}

	rdb, err := newRedis(c.StateRedis)
	if err != nil {
		_ = sqlDB.Close()
		panic(fmt.Errorf("初始化 Redis 失败: %w", err))
	}

	logx.Infof("identity 基础设施初始化完成: postgres=%s:%d redis=%s:%d",
		c.Postgres.Host, c.Postgres.Port, c.StateRedis.Host, c.StateRedis.Port)

	return &ServiceContext{
		Config: c,
		DB:     db,
		SQLDB:  sqlDB,
		Redis:  rdb,
	}
}

func newPostgres(c config.PostgresConf) (*gorm.DB, *sql.DB, error) {
	pg := withPostgresDefaults(c)

	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(pg.User, pg.Password),
		Host:   fmt.Sprintf("%s:%d", pg.Host, pg.Port),
		Path:   pg.DBName,
	}

	query := dsnURL.Query()
	query.Set("sslmode", pg.SSLMode)
	query.Set("TimeZone", pg.TimeZone)
	query.Set("connect_timeout", fmt.Sprintf("%d", pg.ConnectTimeoutSeconds))
	dsnURL.RawQuery = query.Encode()

	db, err := gorm.Open(postgres.Open(dsnURL.String()), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(pg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(pg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(pg.ConnMaxLifetimeSeconds) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(pg.ConnMaxIdleTimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pg.ConnectTimeoutSeconds)*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	return db, sqlDB, nil
}

func newRedis(c config.RedisConf) (*redis.Client, error) {
	rc := withRedisDefaults(c)

	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Username:     rc.Username,
		Password:     rc.Password,
		DB:           rc.DB,
		DialTimeout:  time.Duration(rc.DialTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(rc.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(rc.WriteTimeoutSeconds) * time.Second,
		PoolTimeout:  time.Duration(rc.PoolTimeoutSeconds) * time.Second,
		MaxRetries:   rc.MaxRetries,
		PoolSize:     rc.PoolSize,
		MinIdleConns: rc.MinIdleConns,
	}

	if rc.EnableTLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rc.DialTimeoutSeconds)*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}

func withPostgresDefaults(c config.PostgresConf) config.PostgresConf {
	if c.Port == 0 {
		c.Port = 5432
	}
	if strings.TrimSpace(c.SSLMode) == "" {
		c.SSLMode = "disable"
	}
	if strings.TrimSpace(c.TimeZone) == "" {
		c.TimeZone = "Asia/Shanghai"
	}
	if c.ConnectTimeoutSeconds <= 0 {
		c.ConnectTimeoutSeconds = 5
	}
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = 20
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnMaxLifetimeSeconds <= 0 {
		c.ConnMaxLifetimeSeconds = 1800
	}
	if c.ConnMaxIdleTimeSeconds <= 0 {
		c.ConnMaxIdleTimeSeconds = 600
	}

	return c
}

func withRedisDefaults(c config.RedisConf) config.RedisConf {
	if c.Port == 0 {
		c.Port = 6379
	}
	if c.DialTimeoutSeconds <= 0 {
		c.DialTimeoutSeconds = 3
	}
	if c.ReadTimeoutSeconds <= 0 {
		c.ReadTimeoutSeconds = 3
	}
	if c.WriteTimeoutSeconds <= 0 {
		c.WriteTimeoutSeconds = 3
	}
	if c.PoolTimeoutSeconds <= 0 {
		c.PoolTimeoutSeconds = 4
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}
	if c.PoolSize <= 0 {
		c.PoolSize = 10
	}
	if c.MinIdleConns < 0 {
		c.MinIdleConns = 0
	}

	return c
}
