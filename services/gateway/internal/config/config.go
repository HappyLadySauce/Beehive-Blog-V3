// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"fmt"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// GatewaySecurityConf defines gateway security related settings.
// GatewaySecurityConf 定义网关安全相关配置。
type GatewaySecurityConf struct {
	TokenPrefix         string   `json:"TokenPrefix"`
	TrustedProxyHeaders []string `json:"TrustedProxyHeaders"`
	TrustedProxyCIDRs   []string `json:"TrustedProxyCIDRs"`
}

// InternalRPCConf defines gateway access to an internal RPC service.
// InternalRPCConf 定义 gateway 访问内部 RPC 服务的配置。
type InternalRPCConf struct {
	zrpc.RpcClientConf
	InternalAuthToken  string `json:"InternalAuthToken"`
	InternalCallerName string `json:"InternalCallerName"`
}

// IdentityRPCConf defines gateway access to the identity RPC service.
// IdentityRPCConf 定义 gateway 访问 identity RPC 服务的配置。
type IdentityRPCConf = InternalRPCConf

// ContentRPCConf defines gateway access to the content RPC service.
// ContentRPCConf 定义 gateway 访问 content RPC 服务的配置。
type ContentRPCConf = InternalRPCConf

type Config struct {
	rest.RestConf
	IdentityRPC IdentityRPCConf
	ContentRPC  ContentRPCConf
	Security    GatewaySecurityConf
}

var supportedTrustedProxyHeaders = map[string]struct{}{
	"x-forwarded-for": {},
	"x-real-ip":       {},
	"client-ip":       {},
}

// Validate verifies gateway configuration before startup.
// Validate 在启动前验证 gateway 配置。
func (c Config) Validate() error {
	if _, err := ctxmeta.ParseTrustedProxyCIDRs(c.Security.TrustedProxyCIDRs); err != nil {
		return fmt.Errorf("security.trusted_proxy_cidrs is invalid: %w", err)
	}

	hasHeaders := hasNonEmptyValues(c.Security.TrustedProxyHeaders)
	hasCIDRs := hasNonEmptyValues(c.Security.TrustedProxyCIDRs)
	if hasHeaders && !hasCIDRs {
		return fmt.Errorf("security.trusted_proxy_headers requires non-empty security.trusted_proxy_cidrs")
	}
	if hasCIDRs && !hasHeaders {
		return fmt.Errorf("security.trusted_proxy_cidrs requires non-empty security.trusted_proxy_headers")
	}
	for _, rawHeader := range c.Security.TrustedProxyHeaders {
		header := strings.ToLower(strings.TrimSpace(rawHeader))
		if header == "" {
			continue
		}
		if _, ok := supportedTrustedProxyHeaders[header]; !ok {
			return fmt.Errorf("security.trusted_proxy_headers contains unsupported value %q", rawHeader)
		}
	}
	if strings.TrimSpace(c.IdentityRPC.InternalAuthToken) == "" {
		return fmt.Errorf("IdentityRPC.InternalAuthToken is required")
	}
	if strings.TrimSpace(c.IdentityRPC.InternalCallerName) == "" {
		return fmt.Errorf("IdentityRPC.InternalCallerName is required")
	}
	if strings.TrimSpace(c.ContentRPC.InternalAuthToken) == "" {
		return fmt.Errorf("ContentRPC.InternalAuthToken is required")
	}
	if strings.TrimSpace(c.ContentRPC.InternalCallerName) == "" {
		return fmt.Errorf("ContentRPC.InternalCallerName is required")
	}

	return nil
}

func hasNonEmptyValues(values []string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}
