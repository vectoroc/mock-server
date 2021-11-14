package model

import "encoding/json"

type KeyToValue map[string]string

func (kv *KeyToValue) MarshalJSON() ([]byte, error) {
	if kv == nil || len(*kv) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string]string(*kv))
}

func (kv *KeyToValue) UnmarshalJSON(data []byte) error {
	var list []struct {
		Name  string
		Value string
	}

	*kv = KeyToValue{}

	if err := json.Unmarshal(data, &list); err == nil {
		for _, item := range list {
			(*kv)[item.Name] = item.Value
		}

		return nil
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, val := range m {
		(*kv)[k] = val
	}

	return nil
}
