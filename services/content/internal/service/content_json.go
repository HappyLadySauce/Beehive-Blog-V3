package service

import (
	"encoding/json"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

func stringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func bodyJSONPtr(value string) (*string, error) {
	return jsonPtr(value, "body_json")
}

func metadataJSONPtr(value string) (*string, error) {
	return jsonPtr(value, "metadata_json")
}

func jsonPtr(value, fieldName string) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	if !json.Valid([]byte(trimmed)) {
		return nil, errs.New(errs.CodeContentInvalidArgument, fieldName+" must be valid JSON")
	}
	return &trimmed, nil
}
