package model

import (
	"encoding/json"
	"fmt"
)

// BodyWithContentType - response body
type BodyWithContentType struct {
	Not         bool   `json:"not,omitempty"`
	Type        string `json:"type"`
	Base64Bytes string `json:"base64Bytes,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Json        string `json:"json,omitempty"`
	String      string `json:"string,omitempty"`
	Xml         string `json:"xml,omitempty"`
}

func (b *BodyWithContentType) UnmarshalJSON(data []byte) error {
	var val interface{}
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}

	switch casted := val.(type) {
	case string:
		b.Type = "STRING"
		b.String = casted

	case map[string]interface{}:
		return unmarshalMap(b, casted, data)

	default:
		b.Type = "JSON"
		b.Json = string(data)
	}

	return nil
}

func unmarshalMap(b *BodyWithContentType, data map[string]interface{}, message json.RawMessage) error {
	typeVal, ok := data["type"]
	if !ok {
		b.Type = "JSON"
		b.ContentType = "application/json"
		b.Json = string(message)
		return nil
	}

	if b.Type, ok = typeVal.(string); !ok {
		return fmt.Errorf("type is not a string %w", ErrBadFormat)
	}

	b.ContentType, _ = data["contentType"].(string)
	b.Not, _ = data["not"].(bool)

	switch b.Type {
	case "STRING":
		b.String, ok = data["string"].(string)

	case "JSON":
		b.Json, ok = data["json"].(string)

	case "XML":
		b.Xml, ok = data["xml"].(string)

	case "BINARY":
		b.Base64Bytes, ok = data["base64Bytes"].(string)
	}

	if !ok {
		return fmt.Errorf("unknown body type %s, %w", b.Type, ErrBadFormat)
	}

	return nil
}
