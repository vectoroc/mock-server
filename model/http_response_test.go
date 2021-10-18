package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalHttpResponse(t *testing.T) {
	data := `
{
  "headers": {
    "Location": "https://www.google.com/welcome"
  },
  "cookies": {
    "sessionid": "some-secret"
  },
  "statusCode": 301
}`

	resp := HttpResponse{}
	err := json.Unmarshal([]byte(data), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Headers)
	assert.NotEmpty(t, resp.Cookies)
}
