package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpectations_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshal one expectation with only id", func(t *testing.T) {
		data := `{"id": "1"}`

		v := &Expectations{}
		err := json.Unmarshal([]byte(data), v)
		assert.NoError(t, err)
		require.Len(t, v.list, 1)
		assert.Equal(t, "1", v.list[0].Id)
	})

	t.Run("unmarshal more fields for one expectation", func(t *testing.T) {
		data := `[{
  "httpRequest": {
    "method": "GET",
    "headers": [
      {
        "name": "Host",
        "values": [
          "ya.ru"
        ]
      }
    ]
  },
  "httpBody": {
    "delay": {
      "timeUnit": "SECONDS",
      "value": 10
    },
    "body": {
      "type": "JSON",
      "json": "{\"success\": true}"
    },
    "statusCode": 200
  }
}]`

		v := &Expectations{}
		err := json.Unmarshal([]byte(data), v)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	})
}
