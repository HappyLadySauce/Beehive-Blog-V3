package svc

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	SQLDB    *sql.DB
	Store    *repo.Store
	Services *contentservice.Manager
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("content service context validation failed: %w", err)
	}

	db, sqlDB, err := newPostgres(c.Postgres)
	if err != nil {
		return nil, fmt.Errorf("initialize PostgreSQL failed: %w", err)
	}
	store := repo.NewStore(db)
	readinessChecker := func(ctx context.Context) error {
		if sqlDB == nil {
			return fmt.Errorf("postgres connection is not initialized")
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			return fmt.Errorf("postgres readiness probe failed: %w", err)
		}
		return nil
	}
	services := contentservice.NewManager(contentservice.Dependencies{
		Config:         c,
		Store:          store,
		CheckReadiness: readinessChecker,
	})

	logs.Ctx(context.Background()).Info(
		"content_infrastructure_initialized",
		logs.String("postgres_host", c.Postgres.Host),
		logs.Int("postgres_port", withPostgresDefaults(c.Postgres).Port),
	)

	return &ServiceContext{
		Config:   c,
		DB:       db,
		SQLDB:    sqlDB,
		Store:    store,
		Services: services,
	}, nil
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
