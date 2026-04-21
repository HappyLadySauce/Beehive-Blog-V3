// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"fmt"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// GatewaySecurityConf defines gateway security related settings.
// GatewaySecurityConf 定义网关安全相关配置。
type GatewaySecurityConf struct {
	TokenPrefix         string   `json:",default=Bearer"`
	TrustedProxyHeaders []string `json:",optional"`
	TrustedProxyCIDRs   []string `json:",optional"`
}

type Config struct {
	rest.RestConf
	IdentityRPC zrpc.RpcClientConf
	Security    GatewaySecurityConf
}

// Validate verifies gateway configuration before startup.
// Validate 在启动前验证 gateway 配置。
func (c Config) Validate() error {
	if _, err := ctxmeta.ParseTrustedProxyCIDRs(c.Security.TrustedProxyCIDRs); err != nil {
		return fmt.Errorf("security.trusted_proxy_cidrs is invalid: %w", err)
	}

	return nil
}
