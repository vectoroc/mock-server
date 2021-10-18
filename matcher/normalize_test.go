package matcher

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"testing"
)

func TestNormalizeRequest(t *testing.T) {
	req := model.HttpRequest{
		Headers: model.KeyToMultiValue{
			Values: map[string][]string{
				"host":            {"YA.RU"},
				"content-type":    {"text/plain"},
				"x-forwarded-for": {"1.1.1.1"},
			},
		},
		Method: "get",
		Path:   "robots.txt",
	}

	expectedHeaders := model.KeyToMultiValue{
		Values: map[string][]string{
			"Host":            {"YA.RU"},
			"Content-Type":    {"text/plain"},
			"X-Forwarded-For": {"1.1.1.1"},
		},
	}
	expectedPath := "/robots.txt"
	expcetedMethod := "GET"

	NormalizeRequest(&req)

	assert.EqualValues(t, expectedHeaders, req.Headers)
	assert.Equal(t, expectedPath, req.Path)
	assert.Equal(t, expcetedMethod, req.Method)

	t.Run("normalization should be an idempotent operation", func(t *testing.T) {
		NormalizeRequest(&req)

		assert.EqualValues(t, expectedHeaders, req.Headers)
		assert.Equal(t, expectedPath, req.Path)
		assert.Equal(t, expcetedMethod, req.Method)
	})
}
