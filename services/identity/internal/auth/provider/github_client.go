package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// GitHubClient provides GitHub OAuth start and callback capabilities.
// GitHubClient 提供 GitHub OAuth 的 start 与 callback 能力。
type GitHubClient struct {
	HTTPClient    *http.Client
	OAuthEndpoint oauth2.Endpoint
	APIBaseURL    string
	Conf          config.OAuthProviderConf
}

// NewGitHubClient creates a GitHub OAuth client with production defaults.
// NewGitHubClient 使用生产默认值创建 GitHub OAuth 客户端。
func NewGitHubClient(conf config.OAuthProviderConf) *GitHubClient {
	return &GitHubClient{
		HTTPClient:    &http.Client{Timeout: 10 * time.Second},
		OAuthEndpoint: githuboauth.Endpoint,
		APIBaseURL:    "https://api.github.com",
		Conf:          conf,
	}
}

// Name returns the provider name.
// Name 返回 provider 名称。
func (c *GitHubClient) Name() string {
	return auth.ProviderGitHub
}

// Enabled reports whether the provider is configured and enabled.
// Enabled 判断 provider 是否配置启用。
func (c *GitHubClient) Enabled() bool {
	return c != nil && c.Conf.Enabled
}

// LoginReady reports whether the provider has a complete login implementation.
// LoginReady 判断 provider 是否具备完整登录实现。
func (c *GitHubClient) LoginReady() bool {
	return c != nil &&
		c.Enabled() &&
		strings.TrimSpace(c.Conf.ClientID) != "" &&
		strings.TrimSpace(c.Conf.ClientSecret) != "" &&
		strings.TrimSpace(c.Conf.RedirectURL) != ""
}

// RedirectURL returns the configured redirect URL.
// RedirectURL 返回配置中的回调地址。
func (c *GitHubClient) RedirectURL() string {
	return strings.TrimSpace(c.Conf.RedirectURL)
}

// BuildAuthorizeURL builds the GitHub authorize URL.
// BuildAuthorizeURL 构建 GitHub 授权地址。
func (c *GitHubClient) BuildAuthorizeURL(state string) (string, error) {
	oauthConf := oauth2.Config{
		ClientID:     strings.TrimSpace(c.Conf.ClientID),
		ClientSecret: strings.TrimSpace(c.Conf.ClientSecret),
		RedirectURL:  strings.TrimSpace(c.Conf.RedirectURL),
		Scopes:       trimmedScopes(c.Conf.Scopes),
		Endpoint:     c.OAuthEndpoint,
	}

	return oauthConf.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

// ExchangeCode exchanges a GitHub authorization code for an access token.
// ExchangeCode 使用 GitHub 授权码交换 access token。
func (c *GitHubClient) ExchangeCode(ctx context.Context, code, redirectURI string) (*AccessToken, error) {
	oauthConf := oauth2.Config{
		ClientID:     strings.TrimSpace(c.Conf.ClientID),
		ClientSecret: strings.TrimSpace(c.Conf.ClientSecret),
		RedirectURL:  strings.TrimSpace(c.Conf.RedirectURL),
		Scopes:       trimmedScopes(c.Conf.Scopes),
		Endpoint:     c.OAuthEndpoint,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, c.HTTPClient)
	token, err := oauthConf.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", redirectURI))
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Token: strings.TrimSpace(token.AccessToken),
	}, nil
}

// FetchProfile fetches the authenticated GitHub profile and normalizes it.
// FetchProfile 拉取已授权 GitHub 用户资料并完成标准化。
func (c *GitHubClient) FetchProfile(ctx context.Context, accessToken *AccessToken) (*Profile, []byte, error) {
	if accessToken == nil || strings.TrimSpace(accessToken.Token) == "" {
		return nil, nil, fmt.Errorf("github access token is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.APIBaseURL, "/")+"/user", nil)
	if err != nil {
		return nil, nil, err
	}
	c.applyStandardHeaders(req, accessToken.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("github user API returned status %d", resp.StatusCode)
	}

	var rawMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawMap); err != nil {
		return nil, nil, err
	}

	emailPtr, emailVerified, err := c.fetchVerifiedEmail(ctx, accessToken.Token)
	if err != nil {
		return nil, nil, err
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, nil, err
	}

	var subject string
	if id, ok := rawMap["id"].(float64); ok && int64(id) > 0 {
		subject = fmt.Sprintf("%.0f", id)
	}
	login, _ := rawMap["login"].(string)
	name, _ := rawMap["name"].(string)

	var avatarURLPtr *string
	if avatarURL, ok := rawMap["avatar_url"].(string); ok && strings.TrimSpace(avatarURL) != "" {
		avatarURLPtr = stringPtr(strings.TrimSpace(avatarURL))
	}

	if subject == "" || strings.TrimSpace(login) == "" {
		return nil, nil, fmt.Errorf("github profile is missing required fields")
	}

	return &Profile{
		Subject:          subject,
		SubjectType:      "github_user_id",
		Login:            strings.TrimSpace(login),
		DisplayName:      strings.TrimSpace(name),
		Email:            emailPtr,
		EmailVerified:    emailVerified,
		AvatarURL:        avatarURLPtr,
		RawProfile:       raw,
		ProviderClientID: stringPtr(strings.TrimSpace(c.Conf.ClientID)),
		RequestedScopes:  scopesPtr(trimmedScopes(c.Conf.Scopes)),
	}, raw, nil
}

type githubEmailRecord struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (c *GitHubClient) fetchVerifiedEmail(ctx context.Context, token string) (*string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.APIBaseURL, "/")+"/user/emails", nil)
	if err != nil {
		return nil, false, err
	}
	c.applyStandardHeaders(req, token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(resp.Body)
		return nil, false, nil
	}

	var emails []githubEmailRecord
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, false, nil
	}

	var fallback *string
	for _, item := range emails {
		email := strings.TrimSpace(item.Email)
		if email == "" || !item.Verified {
			continue
		}
		normalized, err := auth.NormalizeEmail(email)
		if err != nil {
			continue
		}
		if item.Primary {
			return stringPtr(normalized), true, nil
		}
		if fallback == nil {
			fallback = stringPtr(normalized)
		}
	}

	if fallback != nil {
		return fallback, true, nil
	}

	return nil, false, nil
}

func (c *GitHubClient) applyStandardHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "beehive-blog-v3-identity")
}
