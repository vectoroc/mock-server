package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyToValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected KeyToValue
	}{
		{
			name:     "object type",
			data:     `{"a": "1", "b": "2"}`,
			expected: KeyToValue{Values: map[string]string{"a": "1", "b": "2"}},
		},
		{
			name:     "list of KV items",
			data:     `[{"name": "a", "value": "1"}, {"name": "b", "value": "3"}]`,
			expected: KeyToValue{Values: map[string]string{"a": "1", "b": "3"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := KeyToValue{}
			err := json.Unmarshal([]byte(test.data), &val)
			assert.NoError(t, err)
		})
	}
}
