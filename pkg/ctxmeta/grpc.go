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

// RequestMeta carries trusted request metadata.
// RequestMeta 承载可信请求元数据。
type RequestMeta struct {
	ForwardedFor string
	RealIP       string
	ClientIP     string
	UserAgent    string
	RequestID    string
}

// GetClientIPFromIncomingContext extracts a trusted client IP from gRPC metadata.
// 从 gRPC metadata 中提取可信的客户端 IP。
func GetClientIPFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	for _, key := range []string{"x-forwarded-for", "x-real-ip", "client-ip"} {
		values := md.Get(key)
		if len(values) == 0 {
			continue
		}

		candidate := strings.TrimSpace(values[0])
		if candidate == "" {
			continue
		}

		if key == "x-forwarded-for" {
			parts := strings.Split(candidate, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}

		return candidate
	}

	return ""
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

// BuildRequestMetaFromHTTP builds trusted metadata from an HTTP request.
// BuildRequestMetaFromHTTP 从 HTTP 请求构建可信元数据。
func BuildRequestMetaFromHTTP(r *http.Request, trustedProxyHeaders []string) RequestMeta {
	forwardedFor := strings.TrimSpace(r.Header.Get(headerXForwardedFor))
	realIP := strings.TrimSpace(r.Header.Get(headerXRealIP))
	clientIP := strings.TrimSpace(r.Header.Get(headerClientIP))

	resolvedClientIP := ExtractTrustedClientIP(r, trustedProxyHeaders)
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

// ExtractTrustedClientIP extracts a trusted client IP by predefined order.
// ExtractTrustedClientIP 按固定顺序提取可信客户端 IP。
func ExtractTrustedClientIP(r *http.Request, trustedProxyHeaders []string) string {
	headerValues := map[string]string{
		headerXForwardedFor: strings.TrimSpace(r.Header.Get(headerXForwardedFor)),
		headerXRealIP:       strings.TrimSpace(r.Header.Get(headerXRealIP)),
		headerClientIP:      strings.TrimSpace(r.Header.Get(headerClientIP)),
	}
	for _, key := range trustedProxyHeaders {
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

	appendIfNotEmpty(headerXForwardedFor, meta.ForwardedFor)
	appendIfNotEmpty(headerXRealIP, meta.RealIP)
	appendIfNotEmpty(headerClientIP, meta.ClientIP)
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
