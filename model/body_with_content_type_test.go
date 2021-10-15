package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBodyWithContentType_UnmarshalJSON(t *testing.T) {
	t.Run("body as a string", func(t *testing.T) {
		data := `"test"`
		res := &BodyWithContentType{}

		err := json.Unmarshal([]byte(data), res)
		assert.NoError(t, err)
		assert.Equal(t, "STRING", res.Type)
		assert.Equal(t, "test", res.String)
	})

	t.Run("body as a raw json", func(t *testing.T) {
		data := `{"foo": [1,2,3], "bar": null}`
		res := &BodyWithContentType{}

		err := json.Unmarshal([]byte(data), res)
		assert.NoError(t, err)
		assert.Equal(t, "JSON", res.Type)
		assert.Equal(t, data, res.Json)
	})

	t.Run("body as a string struct", func(t *testing.T) {
		data := `{"not": true, "type": "STRING", "string": "test str", "contentType": "text/plain"}`
		res := &BodyWithContentType{}

		err := json.Unmarshal([]byte(data), res)
		assert.NoError(t, err)
		assert.EqualValues(t, &BodyWithContentType{
			Not:         true,
			ContentType: "text/plain",
			Type:        "STRING",
			String:      "test str",
		}, res)
	})

	t.Run("body as json with wrong type", func(t *testing.T) {
		data := `{"not": true, "type": "STRING", "json": "true"}`
		res := &BodyWithContentType{}

		err := json.Unmarshal([]byte(data), res)
		assert.Equal(t, ErrBadFormat, err)
	})
}
