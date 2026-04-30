package config

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Postgres          PostgresConf `json:"Postgres"`
	Storage           StorageConf  `json:"Storage"`
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

type StorageConf struct {
	Driver              string           `json:"Driver"`
	PublicBaseURL       string           `json:"PublicBaseURL"`
	PresignTTLSeconds   int              `json:"PresignTTLSeconds"`
	MaxUploadBytes      int64            `json:"MaxUploadBytes"`
	AllowedContentTypes []string         `json:"AllowedContentTypes"`
	Local               LocalStorageConf `json:"Local"`
	S3                  S3StorageConf    `json:"S3"`
}

type LocalStorageConf struct {
	RootDir string `json:"RootDir"`
	TempDir string `json:"TempDir"`
	Bucket  string `json:"Bucket"`
}

type S3StorageConf struct {
	Endpoint        string `json:"Endpoint"`
	Region          string `json:"Region"`
	Bucket          string `json:"Bucket"`
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
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
	if err := c.Storage.Validate(); err != nil {
		return err
	}
	return nil
}

func (c StorageConf) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.Driver)) {
	case "", "local":
		return c.Local.Validate()
	case "s3":
		if strings.TrimSpace(c.PublicBaseURL) == "" {
			return fmt.Errorf("Storage.PublicBaseURL is required for s3 storage")
		}
		return c.S3.Validate()
	default:
		return fmt.Errorf("Storage.Driver must be local or s3")
	}
}

func (c LocalStorageConf) Validate() error {
	if strings.TrimSpace(c.RootDir) == "" {
		return fmt.Errorf("Storage.Local.RootDir is required")
	}
	if strings.TrimSpace(c.TempDir) == "" {
		return fmt.Errorf("Storage.Local.TempDir is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return fmt.Errorf("Storage.Local.Bucket is required")
	}
	return nil
}

func (c S3StorageConf) Validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return fmt.Errorf("Storage.S3.Endpoint is required")
	}
	if strings.TrimSpace(c.Region) == "" {
		return fmt.Errorf("Storage.S3.Region is required")
	}
	if strings.TrimSpace(c.Bucket) == "" {
		return fmt.Errorf("Storage.S3.Bucket is required")
	}
	if strings.TrimSpace(c.AccessKeyID) == "" || strings.TrimSpace(c.SecretAccessKey) == "" {
		return fmt.Errorf("Storage.S3 access credentials are required")
	}
	return nil
}
