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
	headerTrustedIP     = "x-beehive-trusted-client-ip"
	headerUserAgent     = "user-agent"
	headerXRequestID    = "x-request-id"
)

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

// GetClientIPFromIncomingContext extracts a trusted client IP from gRPC metadata.
// 从 gRPC metadata 中提取可信的客户端 IP。
func GetClientIPFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(headerTrustedIP)
	if len(values) == 0 {
		return ""
	}

	return normalizeIP(strings.TrimSpace(values[0]))
}

// GetRequestIDFromIncomingContext extracts request_id from gRPC metadata.
// GetRequestIDFromIncomingContext 从 gRPC metadata 中提取 request_id。
func GetRequestIDFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(headerXRequestID)
	if len(values) == 0 {
		return ""
	}

	return strings.TrimSpace(values[0])
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

// OutgoingContextWithRequestMeta injects metadata into outgoing gRPC context.
// OutgoingContextWithRequestMeta 将元数据注入 gRPC 出站上下文。
func OutgoingContextWithRequestMeta(ctx context.Context, meta RequestMeta) context.Context {
	pairs := []string{}
	appendIfNotEmpty := func(key, value string) {
		if value == "" {
			return
		}
		pairs = append(pairs, key, value)
	}

	appendIfNotEmpty(headerTrustedIP, meta.ClientIP)
	appendIfNotEmpty(headerUserAgent, meta.UserAgent)
	appendIfNotEmpty(headerXRequestID, meta.RequestID)

	if len(pairs) == 0 {
		return ctx
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
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
