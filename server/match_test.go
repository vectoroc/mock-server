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
				"Host": []string{"ya.ru"},
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
				"Host": []string{"ya.ru"},
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
				"Accept-Content": []string{"text/plain"},
			},
		}

		req, err := http.NewRequest("GET", "http://ya.ru/test", nil)
		req.Header.Add("accept-content", "text/plain")
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("match by query params", func(t *testing.T) {
		m := model.HttpRequest{
			QueryStringParameters: map[string][]string{
				"foo": {"bar"},
			},
		}

		req, err := http.NewRequest("GET", "http://ya.ru/test?foo=bar", nil)
		assert.NoError(t, err)

		ok, err := matchHttpRequest(&m, req)
		assert.NoError(t, err)
		assert.True(t, ok)

		t.Run("multiple values match", func(t *testing.T) {
			m := model.HttpRequest{
				QueryStringParameters: map[string][]string{
					"foo": {"bar", "baz"},
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
			name:   "simple text match",
			body:   "some body text",
			expect: model.Body{String: "some body text"},
			match:  true,
		},
		{
			name:   "substring text match",
			body:   "some relatively long text",
			expect: model.Body{SubString: true, String: " relatively "},
			match:  true,
		},
		{
			name:        "text match with content type",
			body:        "text text",
			contentType: "text/plain",
			expect:      model.Body{ContentType: "text/plain", String: "text text"},
			match:       true,
		},
		{
			name:        "text not match",
			body:        "text text 1",
			contentType: "text/plain",
			expect:      model.Body{ContentType: "text/plain", String: "text text 2", Not: true},
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
