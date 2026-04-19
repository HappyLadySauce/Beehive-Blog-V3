package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// GitHubProfile is the subset of GitHub user profile fields we persist.
// GitHubProfile 是我们持久化的 GitHub 用户资料字段子集。
type GitHubProfile struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
}

// ExchangeGitHubCode exchanges an OAuth code for an access token.
// ExchangeGitHubCode 使用 OAuth code 交换 GitHub access token。
func ExchangeGitHubCode(ctx context.Context, conf config.OAuthProviderConf, code, redirectURI string) (string, error) {
	oauthConf := oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
		Endpoint:     githuboauth.Endpoint,
	}

	token, err := oauthConf.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", redirectURI))
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// FetchGitHubProfile fetches the authenticated GitHub profile.
// FetchGitHubProfile 拉取已授权的 GitHub 用户资料。
func FetchGitHubProfile(ctx context.Context, accessToken string) (*GitHubProfile, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "beehive-blog-v3-identity")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("github user API returned status %d", resp.StatusCode)
	}

	var profile GitHubProfile
	raw, err := json.Marshal(map[string]any{})
	if err != nil {
		return nil, nil, err
	}

	var rawMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawMap); err != nil {
		return nil, nil, err
	}
	raw, err = json.Marshal(rawMap)
	if err != nil {
		return nil, nil, err
	}
	if id, ok := rawMap["id"].(float64); ok {
		profile.ID = int64(id)
	}
	if login, ok := rawMap["login"].(string); ok {
		profile.Login = login
	}
	if name, ok := rawMap["name"].(string); ok {
		profile.Name = name
	}
	if email, ok := rawMap["email"].(string); ok && email != "" {
		profile.Email = &email
	}
	if avatarURL, ok := rawMap["avatar_url"].(string); ok && avatarURL != "" {
		profile.AvatarURL = &avatarURL
	}

	if profile.ID == 0 || profile.Login == "" {
		return nil, nil, fmt.Errorf("github profile is missing required fields")
	}

	return &profile, raw, nil
}
