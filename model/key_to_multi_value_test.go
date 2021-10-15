package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyToMultiValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected KeyToMultiValue
	}{
		{
			name:     "object type",
			data:     `{"a": "1", "b": "2"}`,
			expected: KeyToMultiValue{"a": []string{"1"}, "b": []string{"2"}},
		},
		{
			name:     "object with arrays and strings",
			data:     `{"a": ["1", "3"], "b": "2"}`,
			expected: KeyToMultiValue{"a": []string{"1", "3"}, "b": []string{"2"}},
		},
		{
			name:     "list of KV items",
			data:     `[{"name": "a", "values": ["1", "2"]}, {"name": "b", "values": ["3", "4"]}]`,
			expected: KeyToMultiValue{"a": []string{"1", "2"}, "b": []string{"3", "4"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val := KeyToMultiValue{}
			err := json.Unmarshal([]byte(test.data), &val)
			assert.NoError(t, err)
		})
	}
}
