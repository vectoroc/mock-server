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

	t.Run("unmarshal one expectation", func(t *testing.T) {
		data := `{
  "id" : "10c0d440-b54c-4753-a5a8-30a731c2202d",
  "priority" : 0,
  "httpRequest" : {
    "method" : "GET",
    "path" : "/ufwmvutiepksjlbcvalo"
  },
  "httpResponse" : {
    "statusCode" : 200,
    "headers": {
      "Location" : [ "http://dzddrjpjin.livejournal.com/test4" ]
    },
    "body" : "<!DOCTYPE html>\n<html>\n  <head>\n</head><body>\n Hello World!</body>\n</html>"
  },
  "times" : {
    "unlimited" : true
  },
  "timeToLive" : {
    "unlimited" : true
  }
}`

		v := &Expectations{}
		err := json.Unmarshal([]byte(data), v)
		assert.NoError(t, err)
		require.Len(t, v.list, 1)
		require.NotEmpty(t, v.list[0].HttpResponse)
		assert.Equal(t, "10c0d440-b54c-4753-a5a8-30a731c2202d", v.list[0].Id)
	})

	t.Run("unmarshal expectations list", func(t *testing.T) {
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
}, {
  "id" : "12d5c5fd-64d5-4041-af99-de2b7452cc72",
  "priority" : 0,
  "httpRequest" : {
    "method" : "GET",
    "path" : "/test5"
  },
  "httpResponse" : {
    "statusCode" : 302,
    "headers" : {
      "Location" : [ "http://dzddrjpjin.livejournal.com/test4" ]
    }
  },
  "times" : {
    "unlimited" : true
  },
  "timeToLive" : {
    "unlimited" : true
  }
}]`

		v := &Expectations{}
		err := json.Unmarshal([]byte(data), v)
		assert.NoError(t, err)
		assert.NotEmpty(t, v)
	})
}
