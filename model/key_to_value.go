package model

import "encoding/json"

type KeyToValue struct {
	Values map[string]string
}

func (kv *KeyToValue) MarshalJSON() ([]byte, error) {
	if len(kv.Values) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(kv.Values)
}

func (kv *KeyToValue) UnmarshalJSON(data []byte) error {
	var list []struct {
		Name  string
		Value string
	}

	if kv.Values == nil {
		kv.Values = make(map[string]string)
	}

	if err := json.Unmarshal(data, &list); err == nil {
		for _, item := range list {
			kv.Values[item.Name] = item.Value
		}

		return nil
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, val := range m {
		kv.Values[k] = val
	}

	return nil
}
