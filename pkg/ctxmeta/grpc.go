package ctxmeta

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	headerXForwardedFor = "x-forwarded-for"
	headerXRealIP       = "x-real-ip"
	headerClientIP      = "client-ip"
	headerUserAgent     = "user-agent"
	headerXRequestID    = "x-request-id"
)

const (
	// MetadataKeyInternalAuthToken is the metadata key for gateway -> identity shared-secret authentication.
	// MetadataKeyInternalAuthToken 是 gateway -> identity 共享密钥认证的 metadata key。
	MetadataKeyInternalAuthToken = "x-beehive-internal-auth-token"
	// MetadataKeyInternalCaller is the metadata key for the authenticated internal caller name.
	// MetadataKeyInternalCaller 是已认证内部调用方名称的 metadata key。
	MetadataKeyInternalCaller = "x-beehive-internal-caller"
	// MetadataKeyTrustedClientIP is the metadata key for the trusted client IP forwarded by gateway.
	// MetadataKeyTrustedClientIP 是由 gateway 转发的可信客户端 IP metadata key。
	MetadataKeyTrustedClientIP = "x-beehive-trusted-client-ip"
	// MetadataKeyRequestID is the metadata key for the current request identifier.
	// MetadataKeyRequestID 是当前请求标识的 metadata key。
	MetadataKeyRequestID = headerXRequestID
	// MetadataKeyUserAgent is the metadata key for the client user agent.
	// MetadataKeyUserAgent 是客户端 User-Agent 的 metadata key。
	MetadataKeyUserAgent = headerUserAgent
	// MetadataKeyUserID is the metadata key for trusted authenticated user id.
	// MetadataKeyUserID 是可信已认证用户 ID 的 metadata key。
	MetadataKeyUserID = "x-beehive-user-id"
	// MetadataKeySessionID is the metadata key for trusted authenticated session id.
	// MetadataKeySessionID 是可信已认证会话 ID 的 metadata key。
	MetadataKeySessionID = "x-beehive-session-id"
	// MetadataKeyUserRole is the metadata key for trusted authenticated user role.
	// MetadataKeyUserRole 是可信已认证用户角色的 metadata key。
	MetadataKeyUserRole = "x-beehive-user-role"
)

type trustedCallerContextKey struct{}

// RequestMeta carries trusted request metadata.
// RequestMeta 承载可信请求元数据。
type RequestMeta struct {
	ForwardedFor string
	RealIP       string
	ClientIP     string
	UserAgent    string
	RequestID    string
}

// TrustedProxyConfig defines the trusted proxy decision inputs.
// TrustedProxyConfig 定义受信代理判定输入。
type TrustedProxyConfig struct {
	Headers []string
	CIDRs   []*net.IPNet
}

// InternalRPCAuth defines authenticated gateway -> identity RPC metadata.
// InternalRPCAuth 定义已认证的 gateway -> identity RPC metadata。
type InternalRPCAuth struct {
	Token  string
	Caller string
}

// AuthClaims carries trusted authenticated user claims.
// AuthClaims 承载可信已认证用户声明。
type AuthClaims struct {
	UserID    string
	SessionID string
	Role      string
}

// GetClientIPFromIncomingContext extracts a trusted client IP from gRPC metadata.
// 从 gRPC metadata 中提取可信的客户端 IP。
func GetClientIPFromIncomingContext(ctx context.Context) string {
	if _, ok := TrustedInternalCallerFrom(ctx); !ok {
		return ""
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(MetadataKeyTrustedClientIP)
	if len(values) == 0 {
		return ""
	}

	return normalizeIP(strings.TrimSpace(values[0]))
}

// WithTrustedInternalCaller marks the current context as coming from an authenticated internal caller.
// WithTrustedInternalCaller 将当前上下文标记为来自已认证的内部调用方。
func WithTrustedInternalCaller(ctx context.Context, caller string) context.Context {
	return context.WithValue(ctx, trustedCallerContextKey{}, strings.TrimSpace(caller))
}

// TrustedInternalCallerFrom reads the authenticated internal caller from context.
// TrustedInternalCallerFrom 从上下文中读取已认证的内部调用方。
func TrustedInternalCallerFrom(ctx context.Context) (string, bool) {
	caller, ok := ctx.Value(trustedCallerContextKey{}).(string)
	if !ok || strings.TrimSpace(caller) == "" {
		return "", false
	}

	return caller, true
}

// GetRequestIDFromIncomingContext extracts request_id from gRPC metadata.
// GetRequestIDFromIncomingContext 从 gRPC metadata 中提取 request_id。
func GetRequestIDFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(MetadataKeyRequestID)
	if len(values) == 0 {
		return ""
	}

	return strings.TrimSpace(values[0])
}

// TrustedAuthClaimsFromIncomingContext extracts trusted user claims from gRPC metadata.
// TrustedAuthClaimsFromIncomingContext 从 gRPC metadata 中提取可信用户声明。
func TrustedAuthClaimsFromIncomingContext(ctx context.Context) (AuthClaims, bool) {
	if _, ok := TrustedInternalCallerFrom(ctx); !ok {
		return AuthClaims{}, false
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return AuthClaims{}, false
	}

	claims := AuthClaims{
		UserID:    strings.TrimSpace(firstMetadataValue(md, MetadataKeyUserID)),
		SessionID: strings.TrimSpace(firstMetadataValue(md, MetadataKeySessionID)),
		Role:      strings.TrimSpace(firstMetadataValue(md, MetadataKeyUserRole)),
	}
	if claims.UserID == "" {
		return AuthClaims{}, false
	}

	return claims, true
}

// ParseTrustedProxyCIDRs parses trusted proxy CIDR strings once during setup.
// ParseTrustedProxyCIDRs 在初始化阶段一次性解析受信代理 CIDR 字符串。
func ParseTrustedProxyCIDRs(rawCIDRs []string) ([]*net.IPNet, error) {
	if len(rawCIDRs) == 0 {
		return nil, nil
	}

	networks := make([]*net.IPNet, 0, len(rawCIDRs))
	for _, raw := range rawCIDRs {
		cidr := strings.TrimSpace(raw)
		if cidr == "" {
			continue
		}

		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}

	return networks, nil
}

// IsTrustedProxy reports whether the remote address belongs to a trusted proxy.
// IsTrustedProxy 判断远端地址是否属于受信代理。
func IsTrustedProxy(remoteAddr string, cidrs []*net.IPNet) bool {
	if len(cidrs) == 0 {
		return false
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(remoteAddr))
	if err != nil {
		host = strings.TrimSpace(remoteAddr)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	for _, network := range cidrs {
		if network != nil && network.Contains(ip) {
			return true
		}
	}

	return false
}

// BuildRequestMetaFromHTTP builds trusted metadata from an HTTP request.
// BuildRequestMetaFromHTTP 从 HTTP 请求构建可信元数据。
func BuildRequestMetaFromHTTP(r *http.Request, proxyConf TrustedProxyConfig) RequestMeta {
	trustProxyHeaders := IsTrustedProxy(strings.TrimSpace(r.RemoteAddr), proxyConf.CIDRs)

	forwardedFor := ""
	realIP := ""
	clientIP := ""
	if trustProxyHeaders {
		forwardedFor = strings.TrimSpace(r.Header.Get(headerXForwardedFor))
		realIP = strings.TrimSpace(r.Header.Get(headerXRealIP))
		clientIP = strings.TrimSpace(r.Header.Get(headerClientIP))
	}

	resolvedClientIP := ExtractTrustedClientIP(r, proxyConf)
	if clientIP == "" {
		clientIP = resolvedClientIP
	}
	if realIP == "" {
		realIP = resolvedClientIP
	}
	if forwardedFor == "" && resolvedClientIP != "" {
		forwardedFor = resolvedClientIP
	}

	return RequestMeta{
		ForwardedFor: forwardedFor,
		RealIP:       realIP,
		ClientIP:     resolvedClientIP,
		UserAgent:    strings.TrimSpace(r.UserAgent()),
		RequestID:    strings.TrimSpace(r.Header.Get(headerXRequestID)),
	}
}

// ExtractTrustedClientIP extracts a trusted client IP by trusted source and header order.
// ExtractTrustedClientIP 按受信来源和头部顺序提取可信客户端 IP。
func ExtractTrustedClientIP(r *http.Request, proxyConf TrustedProxyConfig) string {
	if IsTrustedProxy(strings.TrimSpace(r.RemoteAddr), proxyConf.CIDRs) {
		headerValues := map[string]string{
			headerXForwardedFor: strings.TrimSpace(r.Header.Get(headerXForwardedFor)),
			headerXRealIP:       strings.TrimSpace(r.Header.Get(headerXRealIP)),
			headerClientIP:      strings.TrimSpace(r.Header.Get(headerClientIP)),
		}
		for _, key := range proxyConf.Headers {
			lowerKey := strings.ToLower(strings.TrimSpace(key))
			rawValue, exists := headerValues[lowerKey]
			if !exists || rawValue == "" {
				continue
			}

			if lowerKey == headerXForwardedFor {
				if firstIP := firstIPFromForwardedFor(rawValue); firstIP != "" {
					return firstIP
				}
				continue
			}

			if normalized := normalizeIP(rawValue); normalized != "" {
				return normalized
			}
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		if normalized := normalizeIP(host); normalized != "" {
			return normalized
		}
	}

	return normalizeIP(strings.TrimSpace(r.RemoteAddr))
}

// BuildIdentityOutgoingContext injects authenticated metadata into gateway -> identity gRPC context.
// BuildIdentityOutgoingContext 将已认证 metadata 注入 gateway -> identity gRPC 上下文。
func BuildIdentityOutgoingContext(ctx context.Context, meta RequestMeta, auth InternalRPCAuth) context.Context {
	return BuildInternalOutgoingContext(ctx, meta, auth, AuthClaims{})
}

// BuildInternalOutgoingContext injects internal auth, request metadata, and optional trusted auth claims.
// BuildInternalOutgoingContext 注入内部认证、请求元数据和可选可信用户声明。
func BuildInternalOutgoingContext(ctx context.Context, meta RequestMeta, auth InternalRPCAuth, claims AuthClaims) context.Context {
	pairs := []string{}
	appendIfNotEmpty := func(key, value string) {
		if value == "" {
			return
		}
		pairs = append(pairs, key, value)
	}

	appendIfNotEmpty(MetadataKeyInternalAuthToken, strings.TrimSpace(auth.Token))
	appendIfNotEmpty(MetadataKeyInternalCaller, strings.TrimSpace(auth.Caller))
	appendIfNotEmpty(MetadataKeyTrustedClientIP, meta.ClientIP)
	appendIfNotEmpty(MetadataKeyUserAgent, meta.UserAgent)
	appendIfNotEmpty(MetadataKeyRequestID, meta.RequestID)
	appendIfNotEmpty(MetadataKeyUserID, claims.UserID)
	appendIfNotEmpty(MetadataKeySessionID, claims.SessionID)
	appendIfNotEmpty(MetadataKeyUserRole, claims.Role)

	if len(pairs) == 0 {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

func firstMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func firstIPFromForwardedFor(raw string) string {
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		if normalized := normalizeIP(strings.TrimSpace(part)); normalized != "" {
			return normalized
		}
	}
	return ""
}

func normalizeIP(raw string) string {
	if raw == "" {
		return ""
	}
	ip := net.ParseIP(raw)
	if ip == nil {
		return ""
	}
	return ip.String()
}
