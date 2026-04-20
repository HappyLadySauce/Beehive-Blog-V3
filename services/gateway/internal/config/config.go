// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// GatewaySecurityConf defines gateway security related settings.
// GatewaySecurityConf 定义网关安全相关配置。
type GatewaySecurityConf struct {
	TokenPrefix         string   `json:",default=Bearer"`
	TrustedProxyHeaders []string `json:",optional"`
}

type Config struct {
	rest.RestConf
	IdentityRPC zrpc.RpcClientConf
	Security    GatewaySecurityConf
}
