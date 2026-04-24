package config

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Postgres          PostgresConf `json:"Postgres"`
	RabbitMQ          RabbitMQConf `json:"RabbitMQ"`
	Outbox            OutboxConf   `json:"Outbox"`
	InternalAuthToken string       `json:"InternalAuthToken"`
	AllowedCallers    []string     `json:"AllowedCallers"`
}

type PostgresConf struct {
	Host                   string `json:"Host"`
	Port                   int    `json:"Port"`
	User                   string `json:"User"`
	Password               string `json:"Password"`
	DBName                 string `json:"DBName"`
	SSLMode                string `json:"SSLMode"`
	TimeZone               string `json:"TimeZone"`
	ConnectTimeoutSeconds  int    `json:"ConnectTimeoutSeconds"`
	MaxOpenConns           int    `json:"MaxOpenConns"`
	MaxIdleConns           int    `json:"MaxIdleConns"`
	ConnMaxLifetimeSeconds int    `json:"ConnMaxLifetimeSeconds"`
	ConnMaxIdleTimeSeconds int    `json:"ConnMaxIdleTimeSeconds"`
}

type RabbitMQConf struct {
	URL                   string `json:"URL"`
	Exchange              string `json:"Exchange"`
	ExchangeType          string `json:"ExchangeType"`
	ConnectTimeoutSeconds int    `json:"ConnectTimeoutSeconds"`
}

type OutboxConf struct {
	DispatchIntervalSeconds int `json:"DispatchIntervalSeconds"`
	BatchSize               int `json:"BatchSize"`
	MaxAttempts             int `json:"MaxAttempts"`
	RetryDelaySeconds       int `json:"RetryDelaySeconds"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Postgres.Host) == "" {
		return fmt.Errorf("Postgres.Host is required")
	}
	if strings.TrimSpace(c.Postgres.User) == "" {
		return fmt.Errorf("Postgres.User is required")
	}
	if strings.TrimSpace(c.Postgres.DBName) == "" {
		return fmt.Errorf("Postgres.DBName is required")
	}
	if strings.TrimSpace(c.RabbitMQ.URL) == "" {
		return fmt.Errorf("RabbitMQ.URL is required")
	}
	if strings.TrimSpace(c.RabbitMQ.Exchange) == "" {
		return fmt.Errorf("RabbitMQ.Exchange is required")
	}
	if c.RabbitMQ.ExchangeType != "" && c.RabbitMQ.ExchangeType != "topic" && c.RabbitMQ.ExchangeType != "direct" && c.RabbitMQ.ExchangeType != "fanout" && c.RabbitMQ.ExchangeType != "headers" {
		return fmt.Errorf("RabbitMQ.ExchangeType is invalid")
	}
	if strings.TrimSpace(c.InternalAuthToken) == "" {
		return fmt.Errorf("InternalAuthToken is required")
	}
	if len(c.AllowedCallers) == 0 {
		return fmt.Errorf("AllowedCallers is required")
	}
	for _, caller := range c.AllowedCallers {
		if strings.TrimSpace(caller) == "" {
			return fmt.Errorf("AllowedCallers must not contain empty values")
		}
	}

	return nil
}
