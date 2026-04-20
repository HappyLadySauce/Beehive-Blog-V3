package svc

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"
	"time"

	identityprovider "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	SQLDB     *sql.DB
	Redis     *redis.Client
	Store     *repo.Store
	Providers *identityprovider.Registry
	Services  *identityservice.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("identity service context validation failed: %w", err))
	}

	db, sqlDB, err := newPostgres(c.Postgres)
	if err != nil {
		panic(fmt.Errorf("initialize PostgreSQL failed: %w", err))
	}

	rdb, err := newRedis(c.StateRedis)
	if err != nil {
		_ = sqlDB.Close()
		panic(fmt.Errorf("initialize Redis failed: %w", err))
	}

	store := repo.NewStore(db)
	providers := identityprovider.NewRegistry(
		identityprovider.NewGitHubClient(c.SSO.GitHub),
		identityprovider.NewQQClient(c.SSO.QQ),
		identityprovider.NewWeChatClient(c.SSO.WeChat),
	)
	services := identityservice.NewManager(c, store, providers)

	logx.Infof("identity infrastructure initialized: postgres=%s:%d redis=%s:%d providers=%d",
		c.Postgres.Host, c.Postgres.Port, c.StateRedis.Host, c.StateRedis.Port, providers.Len())

	return &ServiceContext{
		Config:    c,
		DB:        db,
		SQLDB:     sqlDB,
		Redis:     rdb,
		Store:     store,
		Providers: providers,
		Services:  services,
	}
}

func newPostgres(c config.PostgresConf) (*gorm.DB, *sql.DB, error) {
	pg := withPostgresDefaults(c)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s connect_timeout=%d",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.DBName,
		pg.SSLMode,
		pg.TimeZone,
		pg.ConnectTimeoutSeconds,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
