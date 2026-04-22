package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

const (
	qqAuthorizeURL = "https://graph.qq.com/oauth2.0/authorize"
	qqTokenURL     = "https://graph.qq.com/oauth2.0/token"
	qqOpenIDURL    = "https://graph.qq.com/oauth2.0/me"
	qqUserInfoURL  = "https://graph.qq.com/user/get_user_info"
)

// QQClient provides QQ OAuth start and callback capabilities.
// QQClient 提供 QQ OAuth 的 start 与 callback 能力。
type QQClient struct {
	HTTPClient   *http.Client
	AuthorizeURL string
	TokenURL     string
	OpenIDURL    string
	UserInfoURL  string
	Conf         config.OAuthProviderConf
}

// NewQQClient creates a QQ OAuth client with production defaults.
// NewQQClient 使用生产默认值创建 QQ OAuth 客户端。
func NewQQClient(conf config.OAuthProviderConf) *QQClient {
	return &QQClient{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		AuthorizeURL: qqAuthorizeURL,
		TokenURL:     qqTokenURL,
		OpenIDURL:    qqOpenIDURL,
		UserInfoURL:  qqUserInfoURL,
		Conf:         conf,
	}
}

// Name returns the provider name.
// Name 返回 provider 名称。
func (c *QQClient) Name() string {
	return auth.ProviderQQ
}

// Enabled reports whether the provider is configured and enabled.
// Enabled 判断 provider 是否配置启用。
func (c *QQClient) Enabled() bool {
	return c != nil && c.Conf.Enabled
}

// LoginReady reports whether the provider has a complete login implementation.
// LoginReady 判断 provider 是否具备完整登录实现。
func (c *QQClient) LoginReady() bool {
	return c != nil &&
		c.Enabled() &&
		strings.TrimSpace(c.Conf.ClientID) != "" &&
		strings.TrimSpace(c.Conf.ClientSecret) != "" &&
		strings.TrimSpace(c.Conf.RedirectURL) != ""
}

// RedirectURL returns the configured redirect URL.
// RedirectURL 返回配置中的回调地址。
func (c *QQClient) RedirectURL() string {
	return strings.TrimSpace(c.Conf.RedirectURL)
}

// BuildAuthorizeURL builds the QQ authorize URL.
// BuildAuthorizeURL 构建 QQ 授权地址。
func (c *QQClient) BuildAuthorizeURL(state string) (string, error) {
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", strings.TrimSpace(c.Conf.ClientID))
	query.Set("redirect_uri", strings.TrimSpace(c.Conf.RedirectURL))
	query.Set("state", strings.TrimSpace(state))

	scopes := defaultScopes(c.Conf.Scopes, []string{"get_user_info"})
	if len(scopes) > 0 {
		query.Set("scope", strings.Join(scopes, ","))
	}

	return strings.TrimRight(c.AuthorizeURL, "?") + "?" + query.Encode(), nil
}

// ExchangeCode exchanges a QQ authorization code for an access token.
// ExchangeCode 使用 QQ 授权码交换 access token。
func (c *QQClient) ExchangeCode(ctx context.Context, code, redirectURI string) (*AccessToken, error) {
	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("client_id", strings.TrimSpace(c.Conf.ClientID))
	query.Set("client_secret", strings.TrimSpace(c.Conf.ClientSecret))
	query.Set("code", strings.TrimSpace(code))
	query.Set("redirect_uri", strings.TrimSpace(redirectURI))
	query.Set("fmt", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.TokenURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qq token API returned status %d: %s", resp.StatusCode, bodyPreview(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rawMap, err := decodeJSONObject(body)
	if err != nil {
		return nil, fmt.Errorf("decode qq token response failed: %w (body=%s)", err, bodyPreview(body))
	}
	if code, ok := qqErrorCode(rawMap["error"]); ok && code != "0" {
		return nil, fmt.Errorf("qq token API returned error %s: %s", code, qqErrorMessage(rawMap))
	}

	accessToken := mapString(rawMap, "access_token")
	if strings.TrimSpace(accessToken) == "" {
		return nil, fmt.Errorf("qq token API returned empty access_token (body=%s)", bodyPreview(body))
	}

	openID, err := c.fetchOpenID(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Token:  strings.TrimSpace(accessToken),
		OpenID: openID,
		Scope:  strings.Join(defaultScopes(c.Conf.Scopes, []string{"get_user_info"}), ","),
	}, nil
}

// FetchProfile fetches the authenticated QQ profile and normalizes it.
// FetchProfile 拉取已授权 QQ 用户资料并完成标准化。
func (c *QQClient) FetchProfile(ctx context.Context, accessToken *AccessToken) (*Profile, []byte, error) {
	if accessToken == nil || strings.TrimSpace(accessToken.Token) == "" || strings.TrimSpace(accessToken.OpenID) == "" {
		return nil, nil, fmt.Errorf("qq access token context is incomplete")
	}
	openID := strings.TrimSpace(accessToken.OpenID)

	query := url.Values{}
	query.Set("access_token", strings.TrimSpace(accessToken.Token))
	query.Set("oauth_consumer_key", strings.TrimSpace(c.Conf.ClientID))
	query.Set("openid", openID)
	query.Set("fmt", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.UserInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("qq user info API returned status %d: %s", resp.StatusCode, bodyPreview(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var rawMap map[string]any
	if err := json.Unmarshal(body, &rawMap); err != nil {
		return nil, nil, err
	}
	if ret, ok := qqErrorCode(rawMap["ret"]); ok && ret != "0" {
		return nil, nil, fmt.Errorf("qq user info API returned error %s: %s", ret, bodyPreview(body))
	}

	nickname, _ := rawMap["nickname"].(string)
	avatarURL := firstNonEmptyString(rawMap, "figureurl_qq_2", "figureurl_2", "figureurl_qq_1", "figureurl_1")

	if strings.TrimSpace(openID) == "" || strings.TrimSpace(nickname) == "" {
		return nil, nil, fmt.Errorf("qq profile is missing required fields")
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, nil, err
	}

	return &Profile{
		Subject:          openID,
		SubjectType:      "openid",
		DisplayName:      strings.TrimSpace(nickname),
		AvatarURL:        stringPtr(avatarURL),
		OpenID:           stringPtr(openID),
		RawProfile:       raw,
		ProviderClientID: stringPtr(strings.TrimSpace(c.Conf.ClientID)),
		RequestedScopes:  scopeStringPtr(accessToken.Scope),
	}, raw, nil
}

func (c *QQClient) fetchOpenID(ctx context.Context, accessToken string) (string, error) {
	query := url.Values{}
	query.Set("access_token", strings.TrimSpace(accessToken))
	query.Set("fmt", "json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OpenIDURL+"?"+query.Encode(), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("qq openid API returned status %d: %s", resp.StatusCode, bodyPreview(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	rawMap, err := decodeJSONObject(body)
	if err != nil {
		return "", fmt.Errorf("decode qq openid response failed: %w (body=%s)", err, bodyPreview(body))
	}
	if code, ok := qqErrorCode(rawMap["error"]); ok && code != "0" {
		return "", fmt.Errorf("qq openid API returned error %s: %s", code, qqErrorMessage(rawMap))
	}

	openID := mapString(rawMap, "openid")
	if strings.TrimSpace(openID) == "" {
		return "", fmt.Errorf("qq openid API returned empty openid (body=%s)", bodyPreview(body))
	}

	return strings.TrimSpace(openID), nil
}

func firstNonEmptyString(rawMap map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := rawMap[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func decodeJSONObject(body []byte) (map[string]any, error) {
	var rawMap map[string]any
	if err := json.Unmarshal(body, &rawMap); err != nil {
		return nil, err
	}

	return rawMap, nil
}

func mapString(rawMap map[string]any, key string) string {
	value, ok := rawMap[key]
	if !ok {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return strings.TrimSpace(typed.String())
	case float64:
		if math.Trunc(typed) == typed {
			return fmt.Sprintf("%.0f", typed)
		}
		return fmt.Sprintf("%v", typed)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}

func qqErrorCode(value any) (string, bool) {
	if value == nil {
		return "", false
	}

	switch typed := value.(type) {
	case string:
		trimmed := strings.TrimSpace(typed)
		return trimmed, trimmed != ""
	case json.Number:
		return strings.TrimSpace(typed.String()), true
	case float64:
		if math.Trunc(typed) == typed {
			return fmt.Sprintf("%.0f", typed), true
		}
		return fmt.Sprintf("%v", typed), true
	case int:
		return fmt.Sprintf("%d", typed), true
	case int64:
		return fmt.Sprintf("%d", typed), true
	default:
		text := strings.TrimSpace(fmt.Sprintf("%v", typed))
		return text, text != ""
	}
}

func qqErrorMessage(rawMap map[string]any) string {
	if message := mapString(rawMap, "error_description"); message != "" {
		return message
	}
	if message := mapString(rawMap, "msg"); message != "" {
		return message
	}

	return bodyPreview(mustMarshal(rawMap))
}

func mustMarshal(rawMap map[string]any) []byte {
	body, err := json.Marshal(rawMap)
	if err != nil {
		return nil
	}

	return body
}
