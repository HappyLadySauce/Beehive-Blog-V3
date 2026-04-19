package auth

import "encoding/json"

// MarshalAuditDetail marshals audit detail payloads to JSON.
// MarshalAuditDetail 将审计详情序列化为 JSON。
func MarshalAuditDetail(detail map[string]any) []byte {
	if len(detail) == 0 {
		return nil
	}

	raw, err := json.Marshal(detail)
	if err != nil {
		return nil
	}

	return raw
}
