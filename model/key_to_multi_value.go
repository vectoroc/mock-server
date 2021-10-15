package model

import (
	"encoding/json"
)

type KeyToMultiValue map[string][]string

func (kv KeyToMultiValue) UnmarshalJSON(data []byte) error {
	var list []struct {
		Name   string
		Values []string
	}

	if err := json.Unmarshal(data, &list); err == nil {
		for _, item := range list {
			kv[item.Name] = item.Values
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
			kv[k] = []string{casted}

		case []interface{}:
			kv[k] = make([]string, len(casted))
			var ok bool
			for i, item := range casted {
				kv[k][i], ok = item.(string)
				if !ok {
					return ErrBadFormat
				}
			}

		default:
			return ErrBadFormat
		}
	}

	return nil
}
