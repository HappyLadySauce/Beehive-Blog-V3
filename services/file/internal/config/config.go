package config

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Postgres          PostgresConf      `json:"Postgres"`
	ObjectStorage     ObjectStorageConf `json:"ObjectStorage"`
	InternalAuthToken string            `json:"InternalAuthToken"`
	AllowedCallers    []string          `json:"AllowedCallers"`
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

type ObjectStorageConf struct {
	Enabled                    bool                `json:"Enabled"`
	Endpoint                   string              `json:"Endpoint"`
	Region                     string              `json:"Region"`
	Bucket                     string              `json:"Bucket"`
	AccessKeyID                string              `json:"AccessKeyID"`
	SecretAccessKey            string              `json:"SecretAccessKey"`
	PublicBaseURL              string              `json:"PublicBaseURL"`
	PresignTTLSeconds          int                 `json:"PresignTTLSeconds"`
	MaxBytesByScope            map[string]int64    `json:"MaxBytesByScope"`
	AllowedContentTypesByScope map[string][]string `json:"AllowedContentTypesByScope"`
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
	if c.ObjectStorage.Enabled {
		if err := c.ObjectStorage.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c ObjectStorageConf) Validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return fmt.Errorf("ObjectStorage.Endpoint is required when enabled")
	}
	if strings.TrimSpace(c.Region) == "" {
		return fmt.Errorf("ObjectStorage.Region is required when enabled")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return fmt.Errorf("ObjectStorage.Bucket is required when enabled")
	}
	if strings.TrimSpace(c.AccessKeyID) == "" || strings.TrimSpace(c.SecretAccessKey) == "" {
		return fmt.Errorf("ObjectStorage access credentials are required when enabled")
	}
	if strings.TrimSpace(c.PublicBaseURL) == "" {
		return fmt.Errorf("ObjectStorage.PublicBaseURL is required when enabled")
	}
	return nil
}
