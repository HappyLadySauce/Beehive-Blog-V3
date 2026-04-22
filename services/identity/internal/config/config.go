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
	SSO               SSOConf      `json:"SSO"`
	InternalAuthToken string       `json:"InternalAuthToken"`
	AllowedCallers    []string     `json:"AllowedCallers"`
}

type RedisConf struct {
	Host                string `json:"Host"`
	Port                int    `json:"Port"`
	Username            string `json:"Username"`
	Password            string `json:"Password"`
	DB                  int    `json:"DB"`
	EnableTLS           bool   `json:"EnableTLS"`
	DialTimeoutSeconds  int    `json:"DialTimeoutSeconds"`
	ReadTimeoutSeconds  int    `json:"ReadTimeoutSeconds"`
	WriteTimeoutSeconds int    `json:"WriteTimeoutSeconds"`
	PoolTimeoutSeconds  int    `json:"PoolTimeoutSeconds"`
	MaxRetries          int    `json:"MaxRetries"`
	PoolSize            int    `json:"PoolSize"`
	MinIdleConns        int    `json:"MinIdleConns"`
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

type SecurityConf struct {
	AccessTokenSecret      string `json:"AccessTokenSecret"`
	AccessTokenTTLSeconds  int64  `json:"AccessTokenTTLSeconds"`
	RefreshTokenTTLSeconds int64  `json:"RefreshTokenTTLSeconds"`
	StateTTLSeconds        int64  `json:"StateTTLSeconds"`
	PasswordHashCost       int    `json:"PasswordHashCost"`
}

type SSOConf struct {
	GitHub OAuthProviderConf `json:"GitHub"`
	QQ     OAuthProviderConf `json:"QQ"`
	WeChat OAuthProviderConf `json:"WeChat"`
}

type OAuthProviderConf struct {
	Enabled      bool     `json:"Enabled"`
	ClientID     string   `json:"ClientID"`
	ClientSecret string   `json:"ClientSecret"`
	RedirectURL  string   `json:"RedirectURL"`
	Scopes       []string `json:"Scopes"`
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
	if err := validateOAuthProvider("SSO.QQ", c.SSO.QQ); err != nil {
		return err
	}
	if err := validateOAuthProvider("SSO.WeChat", c.SSO.WeChat); err != nil {
		return err
	}
	if err := validateWeChatScopes(c.SSO.WeChat); err != nil {
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

func validateWeChatScopes(conf OAuthProviderConf) error {
	if !conf.Enabled {
		return nil
	}

	scopes := trimScopes(conf.Scopes)
	if len(scopes) == 0 {
		return nil
	}
	if len(scopes) != 1 || !strings.EqualFold(scopes[0], "snsapi_login") {
		return fmt.Errorf("SSO.WeChat.Scopes only supports snsapi_login")
	}

	return nil
}

func trimScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}

	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		trimmed := strings.TrimSpace(scope)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result
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
