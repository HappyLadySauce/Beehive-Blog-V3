package config

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Postgres          PostgresConf `json:"Postgres"`
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
