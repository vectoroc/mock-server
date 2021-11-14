package model

import (
	"encoding/json"
	"fmt"
)

type KeyToMultiValue map[string][]string

func (kv *KeyToMultiValue) MarshalJSON() ([]byte, error) {
	if kv == nil || len(*kv) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string][]string(*kv))
}

func (kv *KeyToMultiValue) UnmarshalJSON(data []byte) error {
	var list []struct {
		Name   string
		Values []string
	}

	*kv = KeyToMultiValue{}

	if err := json.Unmarshal(data, &list); err == nil {
		for _, item := range list {
			(*kv)[item.Name] = item.Values
		}

		return nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, val := range m {
		switch casted := val.(type) {
		case string:
			(*kv)[k] = []string{casted}

		case []interface{}:
			(*kv)[k] = make([]string, len(casted))
			var ok bool
			for i, item := range casted {
				(*kv)[k][i], ok = item.(string)
				if !ok {
					return fmt.Errorf("expected string, got %T, %w", item, ErrBadFormat)
				}
			}

		default:
			return fmt.Errorf("unexpected type %T, %w", val, ErrBadFormat)
		}
	}

	return nil
}
