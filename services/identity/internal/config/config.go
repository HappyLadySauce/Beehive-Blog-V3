package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Postgres          PostgresConf `json:"Postgres"`
	StateRedis        RedisConf    `json:"StateRedis"`
	Security          SecurityConf `json:"Security"`
	SSO               SSOConf      `json:"SSO,optional"`
	InternalAuthToken string       `json:"InternalAuthToken"`
	AllowedCallers    []string     `json:"AllowedCallers"`
}

type RedisConf struct {
	Host                string `json:"Host"`
	Port                int    `json:"Port,default=6379,range=[1:65536)"`
	Username            string `json:"Username,optional"`
	Password            string `json:"Password,optional"`
	DB                  int    `json:"DB,default=0,range=[0:65536)"`
	EnableTLS           bool   `json:"EnableTLS,optional"`
	DialTimeoutSeconds  int    `json:"DialTimeoutSeconds,default=3,range=[1:61)"`
	ReadTimeoutSeconds  int    `json:"ReadTimeoutSeconds,default=3,range=[1:61)"`
	WriteTimeoutSeconds int    `json:"WriteTimeoutSeconds,default=3,range=[1:61)"`
	PoolTimeoutSeconds  int    `json:"PoolTimeoutSeconds,default=4,range=[1:121)"`
	MaxRetries          int    `json:"MaxRetries,default=1,range=[0:11)"`
	PoolSize            int    `json:"PoolSize,default=10,range=[1:201)"`
	MinIdleConns        int    `json:"MinIdleConns,default=0,range=[0:201)"`
}

type PostgresConf struct {
	Host                   string `json:"Host"`
	Port                   int    `json:"Port,default=5432,range=[1:65536)"`
	User                   string `json:"User"`
	Password               string `json:"Password,optional"`
	DBName                 string `json:"DBName"`
	SSLMode                string `json:"SSLMode,default=disable,options=[disable,require,verify-ca,verify-full]"`
	TimeZone               string `json:"TimeZone,default=Asia/Shanghai"`
	ConnectTimeoutSeconds  int    `json:"ConnectTimeoutSeconds,default=5,range=[1:61)"`
	MaxOpenConns           int    `json:"MaxOpenConns,default=20,range=[1:201)"`
	MaxIdleConns           int    `json:"MaxIdleConns,default=10,range=[0:201)"`
	ConnMaxLifetimeSeconds int    `json:"ConnMaxLifetimeSeconds,default=1800,range=[1:86401)"`
	ConnMaxIdleTimeSeconds int    `json:"ConnMaxIdleTimeSeconds,default=600,range=[1:86401)"`
}

type SecurityConf struct {
	AccessTokenSecret      string `json:"AccessTokenSecret"`
	AccessTokenTTLSeconds  int64  `json:"AccessTokenTTLSeconds,default=900,range=[60:86401)"`
	RefreshTokenTTLSeconds int64  `json:"RefreshTokenTTLSeconds,default=2592000,range=[3600:7776001)"`
	StateTTLSeconds        int64  `json:"StateTTLSeconds,default=600,range=[60:86401)"`
	PasswordHashCost       int    `json:"PasswordHashCost,default=12,range=[10:17)"`
}

type SSOConf struct {
	GitHub OAuthProviderConf `json:"GitHub,optional"`
}

type OAuthProviderConf struct {
	Enabled      bool     `json:"Enabled,optional"`
	ClientID     string   `json:"ClientID,optional"`
	ClientSecret string   `json:"ClientSecret,optional"`
	RedirectURL  string   `json:"RedirectURL,optional"`
	Scopes       []string `json:"Scopes,optional"`
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
	if strings.TrimSpace(c.StateRedis.Host) == "" {
		return fmt.Errorf("StateRedis.Host is required")
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
	if strings.TrimSpace(c.Security.AccessTokenSecret) == "" {
		return fmt.Errorf("Security.AccessTokenSecret is required")
	}
	if c.Security.RefreshTokenTTLSeconds <= c.Security.AccessTokenTTLSeconds {
		return fmt.Errorf("Security.RefreshTokenTTLSeconds must be greater than Security.AccessTokenTTLSeconds")
	}
	if err := validateOAuthProvider("SSO.GitHub", c.SSO.GitHub); err != nil {
		return err
	}

	return nil
}

func validateOAuthProvider(path string, conf OAuthProviderConf) error {
	if !conf.Enabled {
		return nil
	}
	if strings.TrimSpace(conf.ClientID) == "" {
		return fmt.Errorf("%s.ClientID is required", path)
	}
	if strings.TrimSpace(conf.ClientSecret) == "" {
		return fmt.Errorf("%s.ClientSecret is required", path)
	}
	if strings.TrimSpace(conf.RedirectURL) == "" {
		return fmt.Errorf("%s.RedirectURL is required", path)
	}

	return validateOptionalURL(path+".RedirectURL", conf.RedirectURL)
}

func validateOptionalURL(path, raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("%s is not a valid URL: %w", path, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%s must contain a valid scheme and host", path)
	}

	return nil
}
