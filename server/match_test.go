package server

import (
	"github.com/stretchr/testify/assert"
	"mock-server/model"
	"net/http"
	"testing"
)

func TestMatchRequest(t *testing.T) {

	t.Run("match by host", func(t *testing.T) {
		m := model.HttpRequest{
			Headers: model.KeyToMultiValue{
				Values: map[string][]string{
					"Host": {"ya.ru"},
				},
			},
		}

		req, err := http.NewRequest("GET", "http://ya.ru/", nil)
		req.Header.Add("Referer", "http://google.com/")
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("match by host and path", func(t *testing.T) {
		m := model.HttpRequest{
			Headers: model.KeyToMultiValue{
				Values: map[string][]string{
					"Host": []string{"ya.ru"},
				},
			},
			Path: "/test",
		}

		req, err := http.NewRequest("GET", "http://ya.ru/test", nil)
		req.Header.Add("Referer", "http://google.com/")
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)

		t.Run("it should return false for mismatched path", func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://ya.ru/test2", nil)
			assert.NoError(t, err)

			ok, err := matchHttpRequest(&m, req)
			assert.NoError(t, err)
			assert.False(t, ok)
		})
	})

	t.Run("match should normalize header keys", func(t *testing.T) {
		m := model.HttpRequest{
			Headers: model.KeyToMultiValue{
				Values: map[string][]string{
					"Accept-Content": []string{"responseText/plain"},
				},
			},
		}

		req, err := http.NewRequest("GET", "http://ya.ru/test", nil)
		req.Header.Add("accept-content", "responseText/plain")
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("match by query params", func(t *testing.T) {
		m := model.HttpRequest{
			QueryStringParameters: model.KeyToMultiValue{
				Values: map[string][]string{
					"foo": {"bar"},
				},
			},
		}

		req, err := http.NewRequest("GET", "http://ya.ru/test?foo=bar", nil)
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)

		t.Run("multiple values match", func(t *testing.T) {
			m := model.HttpRequest{
				QueryStringParameters: model.KeyToMultiValue{
					Values: map[string][]string{
						"foo": {"bar", "baz"},
					},
				},
			}

			req, err := http.NewRequest("GET", "http://ya.ru/test?foo=bar&foo=baz", nil)
			assert.NoError(t, err)

			ok, err := matchHttpRequest(&m, req)
			assert.NoError(t, err)
			assert.True(t, ok)
		})
	})
}

func TestMatchBody(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		expect      model.Body
		match       bool
	}{
		{
			name:   "simple responseText match",
			body:   "some body responseText",
			expect: model.Body{String: "some body responseText"},
			match:  true,
		},
		{
			name:   "substring responseText match",
			body:   "some relatively long responseText",
			expect: model.Body{SubString: true, String: " relatively "},
			match:  true,
		},
		{
			name:        "responseText match with content type",
			body:        "responseText responseText",
			contentType: "responseText/plain",
			expect:      model.Body{ContentType: "responseText/plain", String: "responseText responseText"},
			match:       true,
		},
		{
			name:        "responseText not match",
			body:        "responseText responseText 1",
			contentType: "responseText/plain",
			expect:      model.Body{ContentType: "responseText/plain", String: "responseText responseText 2", Not: true},
			match:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := MatchBody(test.body, test.contentType, &test.expect)
			assert.Equal(t, test.match, res)
		})
	}
}
